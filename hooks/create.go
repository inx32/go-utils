package hooks

type HookFunc struct {
	Func       func()
	Name       string
	Desc       string
	Concurrent bool
	Weight     uint16
}

type HookNotify struct {
	Chan        chan struct{}
	DoneChan    chan struct{}
	Name        string
	Desc        string
	NonBlocking bool
	Weight      uint16
}
