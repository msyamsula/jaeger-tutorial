package main

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	grpcClient "github.com/msyamsula/jaeger-tutorial/service/a/grpc_client"
	"github.com/msyamsula/jaeger-tutorial/service/a/mytracer"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func f1(ctx context.Context) {
	newCtx, span := mytracer.Tracer.Start(ctx, "f1")
	defer span.End()

	fmt.Println("function 1")

	networkCallHalo(newCtx)
}

func f2(ctx context.Context) {
	_, span := mytracer.Tracer.Start(ctx, "f2")
	defer span.End()

	fmt.Println("function 2")
}

func f3(ctx context.Context) {

	newCtx, span := mytracer.Tracer.Start(ctx, "f3")
	defer span.End()

	f1(newCtx)
	f2(newCtx)
}

func networkCall(ctx context.Context) {

	// use automatic tracing
	// var span trace.Span
	// ctx, span = mytracer.Tracer.Start(ctx, "network call", trace.WithAttributes(semconv.PeerService("service-b")))
	// defer span.End()

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

	// use automatic tracing, don't manually instrument network call
	// var span trace.Span
	// ctx, span = mytracer.Tracer.Start(ctx, "network call halo", trace.WithAttributes(semconv.PeerService("service-b")))
	// defer span.End()

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
	r.Use(otelgin.Middleware(mytracer.Name))
	r.GET("/", func(ctx *gin.Context) {

		// let otelgin do the tracing
		// var span trace.Span
		// var newCtx context.Context
		// newCtx, span := mytracer.Tracer.Start(ctx, "handler-a")
		// defer span.End()

		newCtx := ctx.Request.Context()
		f3(newCtx)

		networkCall(newCtx)

		grpcClient.CallGRPC(newCtx)

		ctx.JSON(http.StatusOK, gin.H{
			"halo": "world",
		})
	})

	r.Run("0.0.0.0:5000")

}
