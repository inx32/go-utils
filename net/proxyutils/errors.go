package proxyutils

import "errors"

var (
	ErrUnsupportedScheme = errors.New("unsupported scheme (only \"socks5\" and \"http\" supported)")
	ErrSocks5Only        = errors.New("only \"socks5\" scheme supported for that")
)
