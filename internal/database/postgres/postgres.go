package postgres

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/exaring/otelpgx"
	pgxzero "github.com/jackc/pgx-zerolog"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rs/zerolog"
	"github.com/shanto-323/axis/config"
	loggerConfig "github.com/shanto-323/axis/pkg/logger"
	"go.opentelemetry.io/otel/trace"
)

type DB struct {
	pool   *pgxpool.Pool
	logger *zerolog.Logger
}

type multiTracer struct {
	tracers []any
}

func (mt *multiTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	for _, tracer := range mt.tracers {
		if t, ok := tracer.(interface {
			TraceQueryStart(context.Context, *pgx.Conn, pgx.TraceQueryStartData) context.Context
		}); ok {
			ctx = t.TraceQueryStart(ctx, conn, data)
		}
	}
	return ctx
}

func (mt *multiTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	for _, tracer := range mt.tracers {
		if t, ok := tracer.(interface {
			TraceQueryEnd(context.Context, *pgx.Conn, pgx.TraceQueryEndData)
		}); ok {
			t.TraceQueryEnd(ctx, conn, data)
		}
	}
}

func New(cfg *config.Config, logger *zerolog.Logger, tracer trace.Tracer) (*DB, error) {
	hostPort := net.JoinHostPort(cfg.Database.Host, strconv.Itoa(cfg.Database.Port))

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		hostPort,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	pgxPoolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx pool cfg: %w", err)
	}

	if tracer != nil {
		pgxPoolConfig.ConnConfig.Tracer = otelpgx.NewTracer()
	}

	logLevel := logger.GetLevel()
	pgxLogger := loggerConfig.NewPgxLogger(logLevel)

	if pgxPoolConfig.ConnConfig.Tracer != nil {
		localTracer := &tracelog.TraceLog{
			Logger:   pgxzero.NewLogger(pgxLogger),
			LogLevel: tracelog.LogLevel(loggerConfig.GetPgxTraceLogLevel(logLevel)),
		}
		pgxPoolConfig.ConnConfig.Tracer = &multiTracer{
			tracers: []any{pgxPoolConfig.ConnConfig.Tracer, localTracer},
		}
	} else {
		pgxPoolConfig.ConnConfig.Tracer = &tracelog.TraceLog{
			Logger:   pgxzero.NewLogger(pgxLogger),
			LogLevel: tracelog.LogLevel(loggerConfig.GetPgxTraceLogLevel(logLevel)),
		}
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxPoolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	logger.Debug().Msg("postgres service initialized successfully")

	return &DB{
		pool:   pool,
		logger: logger,
	}, nil
}

func (db *DB) Ping(ctx context.Context) error {
	return db.pool.Ping(ctx)
}

func (db *DB) IsInitialized(ctx context.Context) bool {
	return db.pool != nil
}

func (db *DB) Close() error {
	db.logger.Info().Msg("closing database connection pool")
	db.pool.Close()
	return nil
}
