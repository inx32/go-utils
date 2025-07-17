package proxyutils

import netproxy "golang.org/x/net/proxy"

type ProxyOpt func(*proxyImpl)

func WithForward(forward netproxy.Dialer) ProxyOpt {
	return func(t *proxyImpl) {
		t.forward = forward
	}
}

func WithAuth(auth ProxyAuth) ProxyOpt {
	return func(t *proxyImpl) {
		auth := &proxyAuthImpl{
			username: auth.Username(),
		}
		if password, ok := auth.Password(); ok {
			auth.password = password
		}
		t.auth = auth
	}
}

type ProxyAuthOpt func(*proxyAuthImpl)

func WithPassword(password string) ProxyAuthOpt {
	return func(t *proxyAuthImpl) {
		t.password = password
	}
}
