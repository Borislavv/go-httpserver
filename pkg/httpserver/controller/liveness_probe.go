package controller

import (
	"context"
	"encoding/json"
	requestctxenum "github.com/Borislavv/go-httpserver/pkg/httpserver/request/ctx"
	"github.com/Borislavv/go-liveness-prober/pkg/liveness"
	"github.com/Borislavv/go-logger/pkg/logger"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"net/http"
)

const K8SProbeGetPath = "/k8s/probe"

type LivenessProbe struct {
	ctx      context.Context
	logger   logger.Logger
	liveness liveness.Prober
}

func NewLivenessProbe(ctx context.Context, logger logger.Logger, liveness liveness.Prober) *LivenessProbe {
	return &LivenessProbe{ctx: ctx, logger: logger, liveness: liveness}
}

func (c *LivenessProbe) Probe(ctx *fasthttp.RequestCtx) {
	reqCtx, ok := ctx.UserValue(requestctxenum.CtxKey).(context.Context)
	if !ok {
		c.logger.ErrorMsg(c.ctx, "context.Context is not exists into the fasthttp.RequestCtx "+
			"(unable to handle request)", nil)
		return
	}

	isAlive := c.liveness.IsAlive()

	resp := make(map[string]map[string]bool, 1)
	resp["data"] = make(map[string]bool, 1)
	resp["data"]["success"] = isAlive

	b, err := json.Marshal(resp)
	if err != nil {
		c.logger.ErrorMsg(reqCtx, "unable to handle request,"+
			" error occurred while marshaling data into []byte", nil)
		return
	}

	if _, err = ctx.Write(b); err != nil {
		c.logger.ErrorMsg(reqCtx, "unable to handle request,"+
			" error occurred while writing data into *fasthttp.RequestCtx", nil)
		return
	}

	if !isAlive {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
	}
}

func (c *LivenessProbe) AddRoute(router *router.Router) {
	router.GET(K8SProbeGetPath, c.Probe)
}
