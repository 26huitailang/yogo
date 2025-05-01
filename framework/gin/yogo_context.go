package gin

import (
	"context"
)

func (ctx *ContainerContext) BaseContext() context.Context {
	return ctx.Request.Context()
}
