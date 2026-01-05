package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/shanto-323/axis/config"
	"github.com/shanto-323/axis/internal/database"
	"github.com/shanto-323/axis/internal/llm"
	"github.com/shanto-323/axis/internal/llm/openrouter"
	"github.com/shanto-323/axis/pkg/tracer"
)

type Server struct {
	Config   *config.Config
	Logger   *zerolog.Logger
	Database database.Database
	LLM      llm.LLM
	Tracer   *tracer.Provider

	httpServer *http.Server
}

func New(cfg *config.Config, logger *zerolog.Logger) (*Server, error) {
	tracer, err := tracer.New(context.Background(), cfg)
	if err != nil {
		return nil, err
	}

	db, err := database.New(cfg, logger, tracer.Tracer)
	if err != nil {
		return nil, err
	}

	llm := openrouter.NewOpenrouter(cfg, logger, tracer.Tracer)

	return &Server{
		Config:   cfg,
		Logger:   logger,
		Database: db,
		LLM:      llm,
		Tracer:   tracer,
	}, nil
}

func (s *Server) SetUpHTTPServer(handler http.Handler) {
	s.httpServer = &http.Server{
		Addr:         ":" + s.Config.Server.Port,
		Handler:      handler,
		ReadTimeout:  time.Duration(s.Config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.Config.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(s.Config.Server.IdleTimeout) * time.Second,
	}
}

func (s *Server) Run() error {
	if s.httpServer == nil {
		return errors.New("http server not initialized")
	}

	s.Logger.Info().
		Str("port", s.Config.Server.Port).
		Str("env", s.Config.Primary.Env).
		Str("db_type", s.Config.Database.Type).
		Msg("starting server")

	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.Database.Close()
}
