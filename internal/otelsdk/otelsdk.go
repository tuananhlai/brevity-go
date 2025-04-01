package otelsdk

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/metric"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/tuananhlai/brevity-go/internal/config"
)

type SetupConfig struct {
	Mode             config.Mode
	ServiceName      string
	CollectorGrpcURL string
}

// Setup bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func Setup(ctx context.Context, cfg SetupConfig) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	resource, err := newResource(cfg.ServiceName)
	if err != nil {
		handleErr(err)
		return
	}

	// Set up logger provider.
	loggerProvider, err := newLoggerProvider(ctx, cfg.CollectorGrpcURL, resource)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)
	global.SetLoggerProvider(loggerProvider)

	tracerProvider, err := newTracerProvider(ctx, cfg.CollectorGrpcURL, resource)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	// Initialize meter provider and application runtime metric producer.
	runtimeProducer := runtime.NewProducer()
	meterProvider, err := newMeterProvider(ctx, cfg.CollectorGrpcURL, resource, runtimeProducer)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second))
	if err != nil {
		handleErr(err)
		return
	}

	return
}

func Tracer(name string) trace.Tracer {
	return otel.Tracer(name)
}

func Meter(name string) metric.Meter {
	return otel.Meter(name)
}

func Logger(name string) *slog.Logger {
	return otelslog.NewLogger(name)
}

func newResource(serviceName string) (*resource.Resource, error) {
	return resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		),
	)
}

func newLoggerProvider(ctx context.Context, endpointURL string, resource *resource.Resource) (*sdklog.LoggerProvider, error) {
	logExporter, err := otlploggrpc.New(ctx, otlploggrpc.WithEndpointURL(endpointURL))
	if err != nil {
		return nil, err
	}

	loggerProvider := sdklog.NewLoggerProvider(
		sdklog.WithResource(resource),
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
	)
	return loggerProvider, nil
}

func newTracerProvider(ctx context.Context, endpointURL string, resource *resource.Resource) (*sdktrace.TracerProvider, error) {
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithEndpointURL(endpointURL))
	if err != nil {
		return nil, err
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(resource),
	)
	return tracerProvider, nil
}

func newMeterProvider(ctx context.Context, endpointURL string, resource *resource.Resource,
	runtimeProducer *runtime.Producer,
) (*sdkmetric.MeterProvider, error) {
	metricExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithEndpointURL(endpointURL))
	if err != nil {
		return nil, err
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(resource),
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(
				metricExporter,
				sdkmetric.WithProducer(runtimeProducer),
				sdkmetric.WithInterval(time.Second*5),
			),
		),
	)
	return meterProvider, nil
}
