package mytracer

import (
	"fmt"
	"io"
	"os"

	"github.com/joho/godotenv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	Name = "service-a"
)

var Tracer trace.Tracer

// NewConsoleExporter returns a console exporter
func NewConsoleExporter(w io.Writer) (sdkTrace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		stdouttrace.WithPrettyPrint(),
		stdouttrace.WithoutTimestamps(),
	)
}

// newExporter returns a jaeger exporter.
func NewExporter(url string) (sdkTrace.SpanExporter, error) {
	eopt := jaeger.WithCollectorEndpoint(
		jaeger.WithEndpoint(url),
	)
	return jaeger.New(eopt)
}

// newResource returns a resource describing this application.
func NewResource(name string) *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(name),
			semconv.ServiceVersion("v0.1.0"),
			attribute.String("environment", "demo"),
		),
	)

	return r
}

func InitTracer() {
	godotenv.Load(".env")

	// JAEGER EXPORTER
	JAEGER_COLLECTOR := os.Getenv("JAEGER_COLLECTOR_URL")
	fmt.Println(JAEGER_COLLECTOR)
	exporter, _ := NewExporter(JAEGER_COLLECTOR)

	// // CONSOLE EXPORTER
	// f, err := os.Create("traces.txt")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// exporter, _ := NewConsoleExporter(f)

	tpa := sdkTrace.NewTracerProvider(
		sdkTrace.WithBatcher(exporter),
		sdkTrace.WithResource(NewResource(Name)),
	)
	otel.SetTracerProvider(tpa)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	Tracer = otel.Tracer(Name)
}
