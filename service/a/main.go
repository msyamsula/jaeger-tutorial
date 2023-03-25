package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	name = "service-a"
)

var tracer trace.Tracer

// newExporter returns a console exporter.
func newExporter(url string) (sdkTrace.SpanExporter, error) {
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
	_, span := tracer.Start(ctx, "f1")
	defer span.End()

	time.Sleep(1 * time.Second)
	fmt.Println("function 1")
}

func f2(ctx context.Context) {
	_, span := tracer.Start(ctx, "f2")
	defer span.End()

	time.Sleep(2 * time.Second)
	fmt.Println("function 2")
}

func f3(ctx context.Context) {

	newCtx, span := tracer.Start(ctx, "f3")
	defer span.End()

	f1(newCtx)
	f2(newCtx)
}

func networkCall(ctx context.Context) {

	_, span := tracer.Start(ctx, "GET /")
	fmt.Println(name)
	defer func() {
		if span.IsRecording() {
			span.SetAttributes(
				attribute.String("http.method", "GET"),
				attribute.String("http.route", "/"),
				attribute.String("kind", "client"),
				attribute.String("net.host.name", "0.0.0.0:5001"),
			)
		}
		span.End()
	}()
	req, _ := http.NewRequestWithContext(ctx, "GET", "http://0.0.0.0:5001", nil)
	client := http.Client{
		// Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	resBody, _ := ioutil.ReadAll(res.Body)
	fmt.Println(resBody)
	res.Body.Close()
}

func initFunc() {
	godotenv.Load(".env")
	JAEGER_COLLECTOR := os.Getenv("JAEGER_COLLECTOR_URL")
	fmt.Println(JAEGER_COLLECTOR)
	exporter, _ := newExporter(JAEGER_COLLECTOR)

	tpa := sdkTrace.NewTracerProvider(
		sdkTrace.WithBatcher(exporter),
		sdkTrace.WithResource(newResource(name)),
	)
	propagator := propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})
	otel.SetTextMapPropagator(propagator)
	otel.SetTracerProvider(tpa)

	tracer = otel.Tracer(name)
}

func main() {

	initFunc()
	r := gin.Default()
	r.Use(otelgin.Middleware(name))
	r.GET("/", func(ctx *gin.Context) {

		// otel.GetTextMapPropagator().Extract(
		// 	ctx, propagation.HeaderCarrier{},
		// )

		// traceState := trace.TraceState{}
		// traceState, _ = traceState.Insert("command", "value")
		// newCtx := trace.ContextWithSpanContext(ctx, trace.NewSpanContext(
		// 	trace.SpanContextConfig{
		// 		TraceState: traceState,
		// 	},
		// ))

		var span trace.Span
		var newCtx context.Context
		newCtx, span = tracer.Start(ctx, "handler-a", trace.WithSpanKind(trace.SpanKindServer))

		defer func() {
			span.End()
		}()

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
