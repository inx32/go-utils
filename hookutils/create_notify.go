package hookutils

type hookNotifyBuilder struct {
	c HookNotify
}

func (b *hookNotifyBuilder) WithChan(c chan struct{})     { b.c.Chan = c }
func (b *hookNotifyBuilder) WithDoneChan(c chan struct{}) { b.c.DoneChan = c }
func (b *hookNotifyBuilder) WithName(name string)         { b.c.Name = name }
func (b *hookNotifyBuilder) WithDesc(desc string)         { b.c.Desc = desc }
func (b *hookNotifyBuilder) WithNonBlocking(flag bool)    { b.c.NonBlocking = flag }
func (b *hookNotifyBuilder) WithWeight(weight uint16)     { b.c.Weight = weight }

func (b *hookNotifyBuilder) Build() *HookNotify { return &b.c }

func NewHookNotifyBuilder() *hookNotifyBuilder {
	return &hookNotifyBuilder{}
}
