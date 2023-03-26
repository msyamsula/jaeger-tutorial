/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a client for Greeter service.
package grpc_client

import (
	"context"
	"fmt"
	"log"

	"github.com/msyamsula/jaeger-tutorial/service/a/mytracer"
	examplePb "github.com/msyamsula/pb-collections/example"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func CallGRPC(ctx context.Context) {

	addr := "0.0.0.0:50051"
	// Set up a connection to the server.
	_, span := mytracer.Tracer.Start(ctx, "init grpc conn")
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	span.End()
	defer conn.Close()

	_, span = mytracer.Tracer.Start(ctx, "init client")
	c := examplePb.NewGreeterClient(conn)
	fmt.Println(c)
	span.End()

	// Contact the server and print out its response.
	additionInput := &examplePb.Numbers{
		A: 3,
		B: 5,
	}
	r, err := c.Addition(ctx, additionInput)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %v", r.GetX())

	multiplicationInput := &examplePb.Numbers{
		A: 2,
		B: 10,
	}

	r, err = c.Multiplication(ctx, multiplicationInput)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %v", r.GetX())
}
