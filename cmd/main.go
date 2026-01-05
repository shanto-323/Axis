package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/shanto-323/axis/config"
	"github.com/shanto-323/axis/internal/database"
	"github.com/shanto-323/axis/internal/server"
	"github.com/shanto-323/axis/internal/server/handler"
	"github.com/shanto-323/axis/internal/server/router"
	"github.com/shanto-323/axis/internal/service"
	logs "github.com/shanto-323/axis/pkg/logger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("error loading  " + err.Error())
	}

	logger, err := logs.New(cfg)
	if err != nil {
		panic("error loading logger " + err.Error())
	}

	if cfg.Database.Type != "mock" {
		if err := database.Migrate(context.Background(), &logger, cfg); err != nil {
			logger.Fatal().Err(err).Msg("failed to migrate database")
		}
	}

	server, err := server.New(cfg, &logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize server")
	}

	services := service.New(server)
	handler := handler.New(server, services)

	router := router.NewRouter(server, handler)
	server.SetUpHTTPServer(router)

	stopChan := make(chan os.Signal, 1)
	errChan := make(chan error, 1)
	signal.Notify(stopChan, os.Interrupt)

	go func() {
		if err := server.Run(); err != nil {
			errChan <- err
		}
	}()

	select {
	case <-stopChan:
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		done := make(chan error, 1)
		go func() {
			done <- server.Stop(ctx)
		}()

		select {
		case <-ctx.Done():
			logger.Fatal().Err(ctx.Err()).Msg("timeout stopping server")
		case err := <-done:
			if err != nil {
				logger.Fatal().Err(err).Msg("error stopping server")
			}
			logger.Info().Msg("server stopped")
		}
	case err := <-errChan:
		logger.Fatal().Err(err).Msg("error running server " + err.Error())
	}
}
