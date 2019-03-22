package api

import (
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type Router struct {
	logger       *zap.Logger
	router       *routing.Router
	ClientID     string
	ClientSecret string
	Host         string
	Port         string
}

// NewRouter create and init new Router
func NewRouter(logger *zap.Logger, id, secret, host, port string) (*Router, error) {
	return &Router{
		logger:       logger,
		router:       routing.New(),
		ClientID:     id,
		ClientSecret: secret,
		Host:         host,
		Port:         port,
	}, nil
}

// RegisterAndRun registers http routers and start serving incoming requests
func (r *Router) RegisterAndRun() error {
	router := r.router

	router.Get("/", r.ShowLoginPage)

	router.Get("/streams", r.ListStreams)
	router.Get("/streams/<name>", r.ShowStreamPage)

	addr := r.Host + ":" + r.Port
	r.logger.Info("HTTP service started", zap.String("URL", "http://"+addr))
	return fasthttp.ListenAndServe(addr, router.HandleRequest)
}
