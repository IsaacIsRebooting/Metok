package launcher

import (
	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// initTracer initializes the TracerProvider for OpenTelemetry.
//
// Parameters:
//
//	serviceName - The name of the service, used to identify the service in the tracing data.
//	endpoint - The endpoint of the tracing data collector.
//
// Return value:
//
//	Returns an error if the TracerProvider fails to initialize.
func initTracer(serviceName string, endpoint string) error {
	// Create a Jaeger exporter instance, specifying the endpoint of the tracing data collector.
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint)))
	if err != nil {
		return err
	}

	// Create a TracerProvider instance, configuring the sampling strategy, batch exporter, and service resource.
	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.ParentBased(trace.TraceIDRatioBased(1.0))),
		trace.WithBatcher(exporter),
		trace.WithResource(
			resource.NewSchemaless(
				semconv.ServiceNameKey.String(serviceName),
				attribute.String("exporter", "jaeger"),
			),
		),
	)

	// Set the global TracerProvider to the one we created.
	otel.SetTracerProvider(tp)

	// Log the initialization information for the TracerProvider.
	log.Infof("Tracing enabled. Exporting spans to %s with service name %s", endpoint, serviceName)

	return nil
}
