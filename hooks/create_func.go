package hooks

type hookFuncBuilder struct {
	c HookFunc
}

func (b *hookFuncBuilder) WithFunc(f func()) { b.c.Func = f }

func (b *hookFuncBuilder) WithName(name string)           { b.c.Name = name }
func (b *hookFuncBuilder) WithDesc(desc string)           { b.c.Desc = desc }
func (b *hookFuncBuilder) WithConcurrent(concurrent bool) { b.c.Concurrent = concurrent }
func (b *hookFuncBuilder) WithWeight(weight uint16)       { b.c.Weight = weight }

func (b *hookFuncBuilder) Build() *HookFunc { return &b.c }

func NewHookFuncBuilder() *hookFuncBuilder {
	return &hookFuncBuilder{}
}
