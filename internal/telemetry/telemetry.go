package telemetry

import (
	"context"
	"log/slog"
	"time"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

// Setup bootstraps the OpenTelemetry pipeline.
func Setup(ctx context.Context) (err error) {
	resource, err := newResource()
	if err != nil {
		return
	}

	// Set up logger provider.
	loggerProvider, err := newLoggerProvider(ctx, resource)
	if err != nil {
		return
	}
	global.SetLoggerProvider(loggerProvider)

	tracerProvider, err := newTracerProvider(ctx, resource)
	if err != nil {
		return
	}
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{}))

	// Initialize meter provider and application runtime metric producer.
	runtimeProducer := runtime.NewProducer()
	meterProvider, err := newMeterProvider(ctx, resource, runtimeProducer)
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

func newResource() (*resource.Resource, error) {
	return resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
		),
	)
}

func newLoggerProvider(ctx context.Context, resource *resource.Resource) (*sdklog.LoggerProvider, error) {
	logExporter, err := otlploghttp.New(ctx)
	if err != nil {
		return nil, err
	}

	loggerProvider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
		sdklog.WithResource(resource),
	)
	return loggerProvider, nil
}

func newTracerProvider(ctx context.Context, resource *resource.Resource) (*sdktrace.TracerProvider, error) {
	traceExporter, err := otlptracehttp.New(ctx)
	if err != nil {
		return nil, err
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(resource),
	)
	return tracerProvider, nil
}

func newMeterProvider(ctx context.Context, resource *resource.Resource,
	runtimeProducer *runtime.Producer,
) (*sdkmetric.MeterProvider, error) {
	metricExporter, err := otlpmetrichttp.New(ctx)
	if err != nil {
		return nil, err
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(resource),
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(
				metricExporter,
				sdkmetric.WithProducer(runtimeProducer),
			),
		),
	)
	return meterProvider, nil
}
