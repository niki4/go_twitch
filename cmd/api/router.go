package api

import (
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

const (
	twitch_api_base = "https://api.twitch.tv/helix/"
)

type Router struct {
	logger *zap.Logger
	router *routing.Router
}

// NewRouter create and init new Router
func NewRouter(logger *zap.Logger) (*Router, error) {
	return &Router{
		logger: logger,
		router: routing.New(),
	}, nil
}

// RegisterAndRun registers http routers and start serving incoming requests
func (r *Router) RegisterAndRun() error {
	router := r.router

	router.Get("/", r.ShowLoginPage)

	router.Get("/streams", r.ListStreams)
	router.Get("/streams/<id>", r.ShowStreamPage)

	r.logger.Info("HTTP service started on:", zap.String("address", "127.0.0.1:6121"))
	return fasthttp.ListenAndServe(":8080", router.HandleRequest)
}
