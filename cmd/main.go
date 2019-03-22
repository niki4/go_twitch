package main

import (
	"github.com/niki4/go_twitch/cmd/api"
	"go.uber.org/zap"
	"os"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	hostName := os.Getenv("HOST")
	if hostName == "" {
		hostName, err = os.Hostname()
		if err != nil {
			logger.Fatal("Host name cannot be defined", zap.Error(err))
		}
	}
	logger.Info("Host name is set", zap.String("Host", hostName))

	listenPort := os.Getenv("PORT")
	if listenPort == "" {
		listenPort = "8080"
	}
	logger.Info("Listening Port is set", zap.String("Port", listenPort))

	clientSecret := os.Getenv("TWITCH_CLIENT_SECRET")
	if clientSecret == "" {
		logger.Fatal("TWITCH_CLIENT_SECRET env variable is not set")
	}
	logger.Info("ClientSecret is set, OK")

	router, err := api.NewRouter(logger, clientSecret, hostName, listenPort)
	if err != nil {
		logger.Fatal("Unable to init NewRouter:", zap.Error(err))
	}

	if err := router.RegisterAndRun(); err != nil {
		logger.Fatal("Router runtime error:", zap.Error(err))
	}
}
