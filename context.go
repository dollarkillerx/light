package light

import "context"

type ContextIF interface {
	context.Context
}

type Context struct {
	ctx context.Context
}
