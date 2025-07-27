package hooks

type HookFunc struct {
	Func func() // non-nil function
	Name string // non-empty string
	Desc string

	// Concurrent determines function concurrency.
	// If it is true, the function will run in a separate goroutine.
	Concurrent bool

	// Weight determines the order in which the functions are started.
	// The functions are run in order of decreasing weight.
	Weight uint16
}

type HookNotify struct {
	Chan chan<- struct{} // non-nil chan

	// DoneChan is a chan that locks Exec until it receives a done event.
	DoneChan <-chan struct{}
	Name     string // non-empty string
	Desc     string

	// NonBlocking determines the chan write lock.
	// If this is true, Exec will not wait for another goroutine to read the chan and
	// will immediately start the next notify.
	//
	// If DoneChan is not nil, Exec will be anyway blocked until the done event is
	// received, even if NonBlocking is true.
	NonBlocking bool

	// Weight determines the order in which the notifies are started.
	// The notifies are run in order of decreasing weight.
	Weight uint16
}
