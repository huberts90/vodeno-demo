package main

import (
	"context"
	"log"
	"os"
	"sync"
	"time"
	"vodeno.com/demo/internal/sender/mock"

	"go.uber.org/zap"

	"vodeno.com/demo/internal/http"
	"vodeno.com/demo/internal/store/postgres"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	httpPort, ok := os.LookupEnv("HTTP_PORT")
	if !ok || httpPort == "" {
		logger.Fatal("HTTP_PORT is not set")
	}

	dbPassword, ok := os.LookupEnv("PG_PASSWORD")
	if !ok || dbPassword == "" {
		logger.Fatal("PG_PASSWORD is not set")
	}

	// TODO: Move parameters to config
	dbParams := &postgres.DatabaseParameters{
		Host:     "postgres",
		Port:     "5432",
		Username: "postgres",
		Password: dbPassword,
		Database: "postgres",
		SSLMode:  "disable",
	}
	postgresDb, err := postgres.Connect(dbParams)
	if err != nil {
		// TODO: Extend log fields with DB params
		logger.Fatal("Failed to connect to PostgreSQL", zap.Error(err))
	}
	err = postgresDb.Ping()
	if err != nil {
		logger.Fatal("Failed to ping postgresql", zap.Error(err))
	}

	postgresStore := postgres.NewSQLStore(postgresDb)

	httpSrv := http.NewServer(logger, postgresStore)

	// TODO: Graceful shutdown
	quit := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		err = httpSrv.Serve(httpPort)
		if err != nil {
			log.Fatal("Failed to create HTTP server", zap.Error(err))
		}
	}()

	// TODO: Make the ticker interval a variable from config
	ticker := time.NewTicker(10 * time.Second)
	sender := mock.NewMocker(logger)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ticker.C:
				sender.Send(context.TODO(), postgresStore)
			case <-quit:
				ticker.Stop()
				return
			}
		}

	}()

	wg.Wait()
	logger.Info("Messenger stopped")
}
