package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/msyamsula/jaeger-tutorial/service/b/mytracer"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func f3(ctx context.Context) {

	_, span := mytracer.Tracer.Start(ctx, "f3")
	defer span.End()

	fmt.Println("from service b, f3")
}

func halo(ctx context.Context) {
	_, span := mytracer.Tracer.Start(ctx, "halo")
	defer span.End()

	fmt.Println("halo")
}

func main() {

	mytracer.InitTracer()

	r := gin.New()
	r.Use(otelgin.Middleware("service-b"))
	r.GET("/", func(ctx *gin.Context) {

		// let otelgin middleware instrument this step
		// newCtx, span := mytracer.Tracer.Start(ctx.Request.Context(), "Incoming request")
		// defer span.End()
		newCtx := ctx.Request.Context()
		f3(newCtx)
		ctx.JSON(http.StatusOK, gin.H{
			"service": "b",
		})
	})

	r.GET("/halo", func(ctx *gin.Context) {

		// let otelgin middleware instrument this step
		// newCtx, span := mytracer.Tracer.Start(ctx.Request.Context(), "Incoming request")
		// defer span.End()

		newCtx := ctx.Request.Context()
		halo(newCtx)
		ctx.JSON(http.StatusOK, gin.H{
			"service": "b halo",
		})
	})

	r.Run("0.0.0.0:5001")
}
