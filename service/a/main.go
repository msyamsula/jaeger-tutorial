package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/msyamsula/jaeger-tutorial/service/a/mytracer"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

func f1(ctx context.Context) {
	newCtx, span := mytracer.Tracer.Start(ctx, "f1")
	defer span.End()

	time.Sleep(1 * time.Second)
	fmt.Println("function 1")

	networkCallHalo(newCtx)
}

func f2(ctx context.Context) {
	_, span := mytracer.Tracer.Start(ctx, "f2")
	defer span.End()

	time.Sleep(2 * time.Second)
	fmt.Println("function 2")
}

func f3(ctx context.Context) {

	newCtx, span := mytracer.Tracer.Start(ctx, "f3")
	defer span.End()

	f1(newCtx)
	f2(newCtx)
}

func networkCall(ctx context.Context) {

	var span trace.Span
	ctx, span = mytracer.Tracer.Start(ctx, "network call", trace.WithAttributes(semconv.PeerService("service-b")))
	defer span.End()

	req, _ := http.NewRequestWithContext(ctx, "GET", "http://0.0.0.0:5001", nil)
	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	resBody, _ := io.ReadAll(res.Body)
	fmt.Println(string(resBody))
	res.Body.Close()
}

func networkCallHalo(ctx context.Context) {

	var span trace.Span
	ctx, span = mytracer.Tracer.Start(ctx, "network call halo", trace.WithAttributes(semconv.PeerService("service-b")))
	defer span.End()

	req, _ := http.NewRequestWithContext(ctx, "GET", "http://0.0.0.0:5001/halo", nil)
	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	resBody, _ := io.ReadAll(res.Body)
	fmt.Println(string(resBody))
	res.Body.Close()
}

func main() {

	mytracer.InitTracer()
	r := gin.Default()
	// r.Use(otelgin.Middleware(mytracer.Name))
	r.GET("/", func(ctx *gin.Context) {

		// var span trace.Span
		// var newCtx context.Context
		newCtx, span := mytracer.Tracer.Start(ctx, "handler-a")
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
