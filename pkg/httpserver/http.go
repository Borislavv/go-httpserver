package httpserver

import (
	"context"
	"errors"
	"github.com/Borislavv/go-httpserver/pkg/httpserver/config"
	"github.com/Borislavv/go-httpserver/pkg/httpserver/controller"
	"github.com/Borislavv/go-httpserver/pkg/httpserver/middleware"
	"github.com/Borislavv/go-logger/pkg/logger"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"sync"
)

type HTTP struct {
	ctx    context.Context
	logger logger.Logger
	server *fasthttp.Server
	config config.Configurator
}

func New(
	ctx context.Context,
	logger logger.Logger,
	controllers []controller.HttpController,
	middlewares []middleware.HttpMiddleware,
) (*HTTP, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	s := &HTTP{ctx: ctx, logger: logger, config: cfg}
	s.initServer(s.buildRouter(controllers), middlewares)
	return s, nil
}

func (s *HTTP) ListenAndServe() {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	wg.Add(1)
	go s.serve(wg)

	wg.Add(1)
	go s.shutdown(wg)
}

func (s *HTTP) serve(wg *sync.WaitGroup) {
	defer wg.Done()

	port := s.config.GetHttpServerPort()

	s.logger.InfoMsg(s.ctx, s.config.GetHttpServerName()+" http server was started", logger.Fields{"port": port})
	defer s.logger.InfoMsg(s.ctx, s.config.GetHttpServerName()+" http server was stopped", logger.Fields{"port": port})

	if err := s.server.ListenAndServe(port); err != nil {
		s.logger.ErrorMsg(s.ctx, err.Error(), logger.Fields{"port": port})
	}
}

func (s *HTTP) shutdown(wg *sync.WaitGroup) {
	defer wg.Done()

	<-s.ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), s.config.GetHttpServerShutDownTimeout())
	defer cancel()

	if err := s.server.ShutdownWithContext(ctx); err != nil {
		if !errors.Is(err, context.Canceled) {
			s.logger.ErrorMsg(s.ctx, err.Error(), logger.Fields{"port": s.config.GetHttpServerPort()})
		}
		return
	}
}

func (s *HTTP) buildRouter(controllers []controller.HttpController) *router.Router {
	r := router.New()
	for _, c := range controllers {
		c.AddRoute(r)
	}
	return r
}

func (s *HTTP) initServer(r *router.Router, mdws []middleware.HttpMiddleware) {
	h := r.Handler

	for i := len(mdws) - 1; i >= 0; i-- {
		h = mdws[i].Middleware(h)
	}

	s.server = &fasthttp.Server{Handler: h}
}
