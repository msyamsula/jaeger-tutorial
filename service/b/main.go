package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"

	"github.com/joho/godotenv"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	name = "service-b"
)

// newExporter returns a console exporter.
func newExporter(url string) (sdkTrace.SpanExporter, error) {

	eopt := jaeger.WithCollectorEndpoint(
		jaeger.WithEndpoint(url),
	)
	return jaeger.New(eopt)
}

// newResource returns a resource describing this application.
func newResource() *resource.Resource {
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

func f3(ctx context.Context) {

	_, span := otel.Tracer(name).Start(ctx, "f3")
	defer span.End()

	time.Sleep(3 * time.Second)
	fmt.Println("from service b, f3")
}

func main() {
	godotenv.Load(".env")
	COLLECTOR_URL := os.Getenv("JAEGER_COLLECTOR_URL")
	fmt.Println(COLLECTOR_URL)
	exporter, err := newExporter(COLLECTOR_URL)
	// f, _ := os.Create("trace.txt")
	// defer f.Close()
	// exporter, err := newFileExporter(f)
	if err != nil {
		log.Fatal(err)
	}

	tp := sdkTrace.NewTracerProvider(
		sdkTrace.WithBatcher(exporter),
		sdkTrace.WithResource(newResource()),
	)
	defer func() {
		if err = tp.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	otel.SetTracerProvider(tp)

	r := gin.Default()
	r.Use(otelgin.Middleware(name))
	r.GET("/", func(ctx *gin.Context) {

		// newCtx := otel.GetTextMapPropagator().Extract(
		// 	ctx, propagation.HeaderCarrier{},
		// )

		// parentSpan := trace.SpanFromContext(ctx)
		// parentBaggage := baggage.FromContext(ctx)

		var span trace.Span
		var newCtx context.Context
		newCtx, span = otel.Tracer(name).Start(ctx, "GET /")
		defer func() {
			if span.IsRecording() {
				span.SetAttributes(
					attribute.String("kind", "server"),
					attribute.String("net.host.name", "0.0.0.0:5001"),
				)
			}
			span.End()
		}()

		f3(newCtx)
		ctx.JSON(http.StatusOK, gin.H{
			"service": "b",
		})
	})

	r.Run("0.0.0.0:5001")

	// f3()
	// fmt.Println("hello")
}
