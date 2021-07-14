package light

import (
	"context"
	"time"
)

type Context struct {
	ctx         context.Context
	metaData    map[string]string
	internalMap map[string]interface{}
}

func NewCtx(ctx context.Context) *Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return &Context{
		ctx:         ctx,
		metaData:    map[string]string{},
		internalMap: map[string]interface{}{},
	}
}

func DefaultCtx() *Context {
	return &Context{
		ctx:         context.Background(),
		metaData:    map[string]string{},
		internalMap: map[string]interface{}{},
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

func (c *Context) SetTimeout(timeout time.Duration) {
	if timeout > 0 {
		c.internalMap["timeout"] = timeout
	}
}

func (c *Context) GetTimeout() time.Duration {
	i, ex := c.internalMap["timeout"]
	if !ex {
		return 0
	}
	return i.(time.Duration)
}

func (c *Context) GetPath() string {
	return c.metaData["light_path"]
}

func (c *Context) SetPath(path string) {
	c.metaData["light_path"] = path
}
