package middleware

import (
	"context"
	"github.com/Borislavv/go-httpserver/pkg/httpserver/config"
	requestctxenum "github.com/Borislavv/go-httpserver/pkg/httpserver/request/ctx"
	"github.com/Borislavv/go-logger/pkg/logger"
	"github.com/savsgio/gotils/uuid"
	"github.com/valyala/fasthttp"
)

const (
	XRequestIDHeader   = "X-Request-ID"
	XRequestGUIDHeader = "X-Request-GUID"
)

type ForwardIDsMiddleware struct {
	ctx    context.Context
	config config.Configurator
	logger logger.Logger
}

func NewForwardIDsMiddleware(ctx context.Context, config config.Configurator, logger logger.Logger) *ForwardIDsMiddleware {
	return &ForwardIDsMiddleware{ctx: ctx, config: config, logger: logger}
}

func (m *ForwardIDsMiddleware) Middleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		// extract x-request-id header from request
		id := string(ctx.Request.Header.Peek(XRequestIDHeader))
		if id == "" {
			id = uuid.V4()
		}

		// extract x-request-guid header from request
		guid := string(ctx.Request.Header.Peek(XRequestGUIDHeader))
		if guid == "" {
			guid = uuid.V4()
		}

		// extract request context
		reqCtx, ok := ctx.UserValue(requestctxenum.CtxKey).(context.Context)
		if !ok {
			m.logger.ErrorMsg(m.ctx, "context.Context is not exists into the fasthttp.RequestCtx "+
				"(unable to forward x-request-id and x-request-guid)", nil)
			next(ctx)
			return
		}

		// build updated context which includes x-request-id and x-request-guid values
		reqCtx = context.WithValue(reqCtx, requestctxenum.ReqID, id)
		reqCtx = context.WithValue(reqCtx, requestctxenum.ReqID, id)

		// set up the updated context into *fasthttp.RequestCtx
		ctx.SetUserValue(requestctxenum.CtxKey, reqCtx)

		// write x-request-id and x-request-guid headers
		ctx.Response.Header.Add(XRequestIDHeader, id)
		ctx.Response.Header.Add(XRequestGUIDHeader, guid)

		next(ctx)
	}
}
