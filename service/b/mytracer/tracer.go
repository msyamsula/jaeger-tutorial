package mytracer

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	Name = "service-b"
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

// NewZipkinExporter returns a zipkin exporter
func NewZipkinExporter(zipkinURL string) (sdkTrace.SpanExporter, error) {
	exporter, err := zipkin.New(zipkinURL)
	if err != nil {
		log.Fatal("zipkin exporter error")
	}

	return exporter, err
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

	// ZIPKIN EXPORTER
	ZIPKIN_COLLECTOR := os.Getenv("ZIPKIN_COLLECTOR_URL")
	fmt.Println(ZIPKIN_COLLECTOR)
	exporter, _ := NewZipkinExporter(ZIPKIN_COLLECTOR)

	// // JAEGER COLLECTOR
	// COLLECTOR_URL := os.Getenv("JAEGER_COLLECTOR_URL")
	// fmt.Println(COLLECTOR_URL)
	// exporter, _ := NewExporter(COLLECTOR_URL)

	// // CONSOLE EXPORTER
	// f, err := os.Create("traces.txt")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// exporter, _ := NewConsoleExporter(f)

	tp := sdkTrace.NewTracerProvider(
		sdkTrace.WithBatcher(exporter),
		sdkTrace.WithResource(NewResource(Name)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	Tracer = otel.Tracer(Name)
	fmt.Println(Tracer)
}
