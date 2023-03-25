package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/msyamsula/jaeger-tutorial/service/b/mytracer"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func f3(ctx context.Context) {

	_, span := mytracer.Tracer.Start(ctx, "f3")
	defer span.End()

	time.Sleep(3 * time.Second)
	fmt.Println("from service b, f3")
}

func halo(ctx context.Context) {
	_, span := mytracer.Tracer.Start(ctx, "halo")
	defer span.End()

	fmt.Println("halo")
}

func main() {

	mytracer.InitTracer()

	// rootHandler := func(w http.ResponseWriter, req *http.Request) {
	// 	ctx := req.Context()
	// 	span := trace.SpanFromContext(ctx)
	// 	defer span.End()
	// 	bag := baggage.FromContext(ctx)
	// 	fmt.Println(bag.Member("username"))

	// 	f3(ctx)

	// 	_, _ = io.WriteString(w, "Hello, world! from service b\n")
	// }

	// otelRootHandler := otelhttp.NewHandler(http.HandlerFunc(rootHandler), "root service b")

	// http.Handle("/", otelRootHandler)
	// err := http.ListenAndServe("0.0.0.0:5001", nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	r := gin.New()
	r.Use(otelgin.Middleware("service-b"))
	r.GET("/", func(ctx *gin.Context) {

		newCtx, span := mytracer.Tracer.Start(ctx.Request.Context(), "Incoming request")
		defer span.End()

		f3(newCtx)
		ctx.JSON(http.StatusOK, gin.H{
			"service": "b",
		})
	})

	r.GET("/halo", func(ctx *gin.Context) {

		newCtx, span := mytracer.Tracer.Start(ctx.Request.Context(), "Incoming request")
		defer span.End()

		halo(newCtx)
		ctx.JSON(http.StatusOK, gin.H{
			"service": "b halo",
		})
	})

	r.Run("0.0.0.0:5001")
}
