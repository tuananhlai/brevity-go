package telemetry

import (
	"context"
	"log/slog"
	"time"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/exporters/autoexport"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// Setup bootstraps the OpenTelemetry pipeline.
func Setup(ctx context.Context) (err error) {
	// Set up logger provider.
	loggerProvider, err := newLoggerProvider(ctx)
	if err != nil {
		return
	}
	global.SetLoggerProvider(loggerProvider)

	tracerProvider, err := newTracerProvider(ctx)
	if err != nil {
		return
	}
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	meterProvider, err := newMeterProvider(ctx)
	if err != nil {
		return
	}
	otel.SetMeterProvider(meterProvider)

	err = runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second))
	if err != nil {
		return
	}

	return
}

// Tracer returns a named OpenTelemetry tracer. It should only
// be called after `Setup` has been called.
func Tracer(name string) trace.Tracer {
	return otel.Tracer(name)
}

// Meter returns a named OpenTelemetry meter. It should only
// be called after `Setup` has been called.
func Meter(name string) metric.Meter {
	return otel.Meter(name)
}

// Logger returns a named OpenTelemetry logger. It should only
// be called after `Setup` has been called.
func Logger(name string) *slog.Logger {
	return otelslog.NewLogger(name)
}

func newLoggerProvider(ctx context.Context) (*sdklog.LoggerProvider, error) {
	logExporter, err := autoexport.NewLogExporter(ctx)
	if err != nil {
		return nil, err
	}

	loggerProvider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
	)
	return loggerProvider, nil
}

func newTracerProvider(ctx context.Context) (*sdktrace.TracerProvider, error) {
	traceExporter, err := autoexport.NewSpanExporter(ctx)
	if err != nil {
		return nil, err
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
	)
	return tracerProvider, nil
}

func newMeterProvider(ctx context.Context) (*sdkmetric.MeterProvider, error) {
	metricReader, err := autoexport.NewMetricReader(ctx)
	if err != nil {
		return nil, err
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(metricReader),
	)
	return meterProvider, nil
}
