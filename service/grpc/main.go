package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/msyamsula/jaeger-tutorial/service/grpc/mytracer"
	pb "github.com/msyamsula/pb-collections/example"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

var (
	port = 50051
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello Handler
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	_, span := mytracer.Tracer.Start(ctx, "say hello")
	defer span.End()
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

// Addition Handler
func (s *server) Addition(ctx context.Context, in *pb.Numbers) (result *pb.Result, err error) {
	_, span := mytracer.Tracer.Start(ctx, "addition")
	defer span.End()
	result = &pb.Result{}
	result.X = in.A + in.B
	return result, err
}

// Multiplication Handler
func (s *server) Multiplication(ctx context.Context, in *pb.Numbers) (result *pb.Result, err error) {
	_, span := mytracer.Tracer.Start(ctx, "multiplication")
	defer span.End()

	result = &pb.Result{}
	result.X = in.A * in.B
	return result, err
}

// run the server
func main() {

	mytracer.InitTracer()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
	)
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
