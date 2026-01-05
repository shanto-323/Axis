package tracer

import (
	"context"
	"fmt"

	"github.com/shanto-323/axis/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
)

type Provider struct {
	Tracer trace.Tracer
}

func New(ctx context.Context, cfg *config.Config) (*Provider, error) {
	exporter, err := otlptracehttp.New(context.Background(),
		otlptracehttp.WithEndpoint("tempo:4318"),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter %w", err)
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(cfg.Observability.ServiceName),
			semconv.DeploymentEnvironmentName(cfg.Observability.Environment),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),

		// !NOTE this is only for development process
		// it will push single trace one at a time no bulk
		sdktrace.WithSyncer(exporter),
	)

	otel.SetTracerProvider(tp)
	tracer := tp.Tracer(cfg.Observability.ServiceName)

	return &Provider{
		Tracer: tracer,
	}, nil
}
