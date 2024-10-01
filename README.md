## HTTP SERVER: Wrapper over fasthttp.

### Usage: 

    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

	cfg, err := sharedconfig.Load()
	if err != nil {
		return nil, err
	}

	lgr, lgrCancel, err := logger.NewLogrus(output)
	if err != nil {
		return nil, err
	}
    defer lgrCancel()

	sharedserver.
		NewHTTP(
			ctx,
			lgr,
			cfg,
			[]controller.HttpController{
                        controller.NewK8SProbe(ctx, logger, liveness),
                    },
			[]middleware.HttpMiddleware{
                        /** exec 1st. */ middleware.NewInitCtxMiddleware(ctx, config),
                        /** exec 2nd. */ middleware.NewApplicationJsonMiddleware(),
                    },
		).
            ListenAndServe()

### ENV: 

    type Config struct {
        // HttpServerName is a name of the shared server.
        HttpServerName string `envconfig:"HTTP_SERVER_NAME" default:"http_server"`
        // HttpServerPort is a port for shared server (endpoints like a /probe for k8s).
        HttpServerPort string `envconfig:"HTTP_SERVER_PORT" default:":8000"`
        // HttpServerShutDownTimeout is a duration value before the server will be closed forcefully.
        HttpServerShutDownTimeout time.Duration `envconfig:"HTTP_SERVER_SHUTDOWN_TIMEOUT" default:"5s"`
        // HttpServerRequestTimeout is a timeout value for close request forcefully.
        HttpServerRequestTimeout time.Duration `envconfig:"HTTP_SERVER_REQUEST_TIMEOUT" default:"1m"`
    }
