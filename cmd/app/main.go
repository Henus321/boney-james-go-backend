package main

import (
	"context"
	"fmt"
	"github.com/Henus321/boney-james-go-backend/internal/config"
	"github.com/Henus321/boney-james-go-backend/pkg/client/mongodb"
	"github.com/Henus321/boney-james-go-backend/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	"time"
)

func main() {
	logging.Init()
	logger := logging.GetLogger()

	logger.Info("router init")
	router := httprouter.New()

	logger.Info("config init")
	cfg := config.GetConfig()

	logger.Info("db init")
	_, err := mongodb.NewClient(context.TODO(), cfg.Storage)
	if err != nil {
		logger.Fatalf("%v", err)
	}

	start(router, cfg)
}

func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("start application")

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.Host, cfg.Listen.Port))

	if err != nil {
		logger.Fatalf("Error, %v", err)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Infof("server running on http://%v:%v", cfg.Listen.Host, cfg.Listen.Port)
	logger.Fatal(server.Serve(listener))
}
