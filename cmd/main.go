package main

import (
	"github.com/niki4/go_twitch/cmd/api"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	router, err := api.NewRouter(logger)
	if err != nil {
		logger.Fatal("Unable to init NewRouter:", zap.Error(err))
	}

	if err := router.RegisterAndRun(); err != nil {
		logger.Fatal("Router runtime error:", zap.Error(err))
	}
}
