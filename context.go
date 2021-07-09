package light

import "context"

type Context struct {
	ctx      context.Context
	metaData map[string]string
}

func NewCtx(ctx context.Context) *Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return &Context{
		ctx:      ctx,
		metaData: map[string]string{},
	}
}

func DefaultCtx() *Context {
	return &Context{
		ctx:      context.Background(),
		metaData: map[string]string{},
	}
}

func (c *Context) Value(key string) string {
	return c.metaData[key]
}

func (c *Context) SetValue(key, val string) {
	c.metaData[key] = val
}

func (c *Context) GetMetaData() map[string]string {
	return c.metaData
}

func (c *Context) SetMetaData(r map[string]string) {
	c.metaData = r
}
