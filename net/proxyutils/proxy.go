package proxyutils

import (
	"context"
	"net"
	"net/http"
	neturl "net/url"
	"os"
	"strconv"

	netproxy "golang.org/x/net/proxy"
)

type Proxy interface {
	// URL returns proxy URL as [net/url.URL].
	URL() *neturl.URL

	// String returns proxy URL as string.
	String() string

	// Dialer returns a dialer what implements the [golang.org/x/net/proxy.Dialer] interface.
	// Only SOCKS5 proxies supported.
	Dialer() (netproxy.Dialer, error)

	// Dial opens a new network connection.
	// Only SOCKS5 proxies supported.
	Dial(network, addr string) (net.Conn, error)

	// DialContext opens a new network connection and closes it when context closes.
	// Only SOCKS5 proxies supported.
	DialContext(ctx context.Context, network, addr string) (net.Conn, error)

	// Transport returns a new [net/http.Transport] what sends HTTP request through the proxy.
	Transport() *http.Transport

	// ModifyTransport modifies existing [net/http.Transport] to make it use the proxy.
	ModifyTransport(transport *http.Transport)

	// RoundTrip sends an HTTP request through proxy.
	RoundTrip(req *http.Request) (*http.Response, error)
}

// implementation of [Proxy]
type proxyImpl struct {
	scheme          string
	hostport        string
	auth            *proxyAuthImpl
	forward         netproxy.Dialer
	cachedtransport *http.Transport
	cacheddialer    netproxy.Dialer
}

func (t *proxyImpl) checkScheme() error {
	if t.scheme != "socks5" && t.scheme != "http" {
		return ErrUnsupportedScheme
	}
	return nil
}

func (t *proxyImpl) checkPort() error {
	_, port, err := net.SplitHostPort(t.hostport)
	if err != nil {
		return err
	}
	if _, err = strconv.ParseUint(port, 10, 16); err != nil {
		return err
	}
	return nil
}

func (t *proxyImpl) check() error {
	if err := t.checkScheme(); err != nil {
		return err
	}
	if err := t.checkPort(); err != nil {
		return err
	}
	return nil
}

func (t *proxyImpl) URL() *neturl.URL {
	var u *neturl.Userinfo

	if t.auth != nil {
		if t.auth.password == "" {
			u = neturl.User(t.auth.username)
		} else {
			u = neturl.UserPassword(t.auth.username, t.auth.password)
		}
	}

	return &neturl.URL{
		Scheme: t.scheme,
		Host:   t.hostport,
		User:   u,
	}
}

func (t *proxyImpl) String() string {
	return t.URL().String()
}

func (t *proxyImpl) Dialer() (netproxy.Dialer, error) {
	if t.cacheddialer != nil {
		return t.cacheddialer, nil
	}
	if t.scheme != "socks5" {
		return nil, ErrSocks5Only
	}

	dialer, err := netproxy.FromURL(t.URL(), t.forward)
	if err != nil {
		return nil, err
	}

	t.cacheddialer = dialer
	return dialer, nil
}

func (t *proxyImpl) Dial(network, addr string) (net.Conn, error) {
	dialer, err := t.Dialer()
	if err != nil {
		return nil, err
	}

	return dialer.Dial(network, addr)
}

func (t *proxyImpl) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	conn, err := t.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		if conn != nil {
			conn.Close()
		}
	}()

	return conn, nil
}

func (t *proxyImpl) Transport() *http.Transport {
	if t.cachedtransport != nil {
		return t.cachedtransport
	}
	transport := &http.Transport{}
	t.ModifyTransport(transport)
	t.cachedtransport = transport
	return transport
}

func (t *proxyImpl) ModifyTransport(transport *http.Transport) {
	switch t.scheme {
	case "socks5":
		transport.DialContext = t.DialContext
	case "http":
		transport.Proxy = http.ProxyURL(t.URL())
	}
}

func (t *proxyImpl) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.Transport().RoundTrip(req)
}

//

func New(scheme, hostport string, opts ...ProxyOpt) (Proxy, error) {
	proxy := &proxyImpl{}
	for _, opt := range opts {
		opt(proxy)
	}
	if err := proxy.check(); err != nil {
		return nil, err
	}
	return proxy, nil
}

func FromURL(url string, opts ...ProxyOpt) (Proxy, error) {
	u, err := neturl.Parse(url)
	if err != nil {
		return nil, err
	}
	return New(u.Scheme, u.Host, opts...)
}

func FromEnvironment(opts ...ProxyOpt) (Proxy, error) {
	url := os.Getenv("ALL_PROXY")
	if url == "" {
		url = os.Getenv("all_proxy")
	}
	if url == "" {
		return nil, nil
	}
	return FromURL(url, opts...)
}
