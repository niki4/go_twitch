package api

import (
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type Router struct {
	logger       *zap.Logger
	router       *routing.Router
	ClientSecret string
	Port         string
}

// NewRouter create and init new Router
func NewRouter(logger *zap.Logger, secret, port string) (*Router, error) {
	return &Router{
		logger:       logger,
		router:       routing.New(),
		ClientSecret: secret,
		Port:         port,
	}, nil
}

// RegisterAndRun registers http routers and start serving incoming requests
func (r *Router) RegisterAndRun() error {
	router := r.router

	router.Get("/", r.ShowLoginPage)

	router.Get("/streams", r.ListStreams)
	router.Get("/streams/<name>", r.ShowStreamPage)

	r.logger.Info("HTTP service started")
	return fasthttp.ListenAndServe(":"+r.Port, router.HandleRequest)
}
