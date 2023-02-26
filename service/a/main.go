package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

const (
	name = "service-a"
)

// newExporter returns a console exporter.
func newExporter(url string) (trace.SpanExporter, error) {

	eopt := jaeger.WithCollectorEndpoint(
		jaeger.WithEndpoint(url),
	)
	return jaeger.New(eopt)
}

// newResource returns a resource describing this application.
func newResource(name string) *resource.Resource {
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

func f1(ctx context.Context) {
	_, span := otel.Tracer(name).Start(ctx, "f1")
	defer span.End()

	// time.Sleep(1 * time.Second)
	fmt.Println("function 1")
}

func f2(ctx context.Context) {
	_, span := otel.Tracer(name).Start(ctx, "f2")
	defer span.End()

	// time.Sleep(2 * time.Second)
	fmt.Println("function 2")
}

func f3(ctx context.Context) {

	newCtx, span := otel.Tracer(name).Start(ctx, "f3")
	defer span.End()

	f1(newCtx)
	f2(newCtx)
}

func networkCall(ctx context.Context) {

	_, span := otel.Tracer(name).Start(ctx, "networkCall")
	defer span.End()

	req, _ := http.NewRequest("GET", "http://0.0.0.0:5001", nil)
	req = req.WithContext(ctx)
	res, _ := http.DefaultClient.Do(req)
	resBody, _ := ioutil.ReadAll(res.Body)
	fmt.Println(resBody)
	res.Body.Close()
}

func main() {
	godotenv.Load(".env")
	JAEGER_COLLECTOR := os.Getenv("JAEGER_COLLECTOR_URL")
	fmt.Println(JAEGER_COLLECTOR)
	exporter, _ := newExporter(JAEGER_COLLECTOR)

	tpa := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(newResource("service-a")),
	)

	otel.SetTracerProvider(tpa)

	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {

		newCtx, span := otel.Tracer(name).Start(ctx, "handler-a")
		defer span.End()

		f3(newCtx)

		networkCall(newCtx)

		ctx.JSON(http.StatusOK, gin.H{
			"halo": "world",
		})
	})

	r.Run("0.0.0.0:5000")

	// f3()
	// fmt.Println("hello")
}
