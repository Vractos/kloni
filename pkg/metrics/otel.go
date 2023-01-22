package metrics

import (
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

type OpenTel struct {
	ServiceName      string
	ServiceVersion   string
	ExporterEndpoint string
	Environment      string
}

func NewOpenTel() *OpenTel {
	return &OpenTel{}
}

func (o *OpenTel) GetTracer() trace.Tracer {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(o.ExporterEndpoint)))
	if err != nil {
		log.Fatal(err)
	}

	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(o.ServiceName),
			semconv.ServiceVersionKey.String(o.ServiceVersion),
			attribute.String("environment", o.Environment),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(r),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	tracer := otel.Tracer("dolly")
	return tracer
}
