package middleware

import (
	"context"
	"github.com/Borislavv/go-httpserver/pkg/httpserver/config"
	requestctxenum "github.com/Borislavv/go-httpserver/pkg/httpserver/request/ctx"
	"github.com/valyala/fasthttp"
)

type InitCtxMiddleware struct {
	ctx    context.Context
	config config.Configurator
}

func NewInitCtxMiddleware(ctx context.Context, config config.Configurator) *InitCtxMiddleware {
	return &InitCtxMiddleware{ctx: ctx, config: config}
}

func (m *InitCtxMiddleware) Middleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		reqCtx, reqCtxCancel := context.WithTimeout(m.ctx, m.config.GetHttpServerRequestTimeout())
		_ = reqCtxCancel

		ctx.SetUserValue(requestctxenum.CtxKey, reqCtx)
		ctx.SetUserValue(requestctxenum.CtxCancelKey, reqCtxCancel)

		next(ctx)
	}
}
