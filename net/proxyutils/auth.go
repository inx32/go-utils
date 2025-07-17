package proxyutils

type ProxyAuth interface {
	Username() string
	Password() (string, bool)
}

type proxyAuthImpl struct {
	username string
	password string
}

func (t *proxyAuthImpl) Username() string {
	return t.username
}

func (t *proxyAuthImpl) Password() (string, bool) {
	if t.password == "" {
		return "", false
	}
	return t.password, true
}

func NewProxyAuth(username string, opts ...ProxyAuthOpt) ProxyAuth {
	auth := &proxyAuthImpl{}
	for _, opt := range opts {
		opt(auth)
	}
	return auth
}
