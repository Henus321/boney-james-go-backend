package main

import (
	"context"
	"fmt"
	"github.com/Henus321/boney-james-go-backend/internal/config"
	"github.com/Henus321/boney-james-go-backend/internal/service/coat"
	"github.com/Henus321/boney-james-go-backend/internal/service/shop"
	"github.com/Henus321/boney-james-go-backend/pkg/client/postgresql"
	"github.com/Henus321/boney-james-go-backend/pkg/logging"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	"time"
)

func main() {
	logging.Init()
	logger := logging.GetLogger()

	router := httprouter.New()

	cfg := config.GetConfig()

	// ??? ошибки на этом уровне обрабатывать или внутри ок?
	db := initDatabase(cfg, logger)

	initCoatService(db, router, logger)

	initShopService(db, router, logger)

	initServer(router, cfg)
}

func initServer(router *httprouter.Router, cfg *config.Config) {
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

func initDatabase(cfg *config.Config, logger *logging.Logger) *pgxpool.Pool {
	logger.Info("db init")
	postgreSQLClient, err := postgresql.NewClient(context.TODO(), 3, cfg.Storage)
	if err != nil {
		logger.Fatalf("%v", err)
	}

	return postgreSQLClient
}

func initCoatService(db *pgxpool.Pool, router *httprouter.Router, logger *logging.Logger) {
	coatStorage := coat.NewStorage(db, logger)
	coatService := coat.NewService(coatStorage)
	coatHandler := coat.NewHandler(coatService, logger)
	logger.Info("coat handler register")
	coatHandler.Register(router)
}

func initShopService(db *pgxpool.Pool, router *httprouter.Router, logger *logging.Logger) {
	shopStorage := shop.NewStorage(db, logger)
	shopService := shop.NewService(shopStorage)
	shopHandler := shop.NewHandler(shopService, logger)
	logger.Info("shop handler register")
	shopHandler.Register(router)
}
