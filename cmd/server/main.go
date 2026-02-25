package main

import (
	"context"
	"fmt"
	"log"
	"net"

	authv1 "github.com/97vaibhav/todo-crud-wire-tutorial-grpc/gen/auth/v1"
	todov1 "github.com/97vaibhav/todo-crud-wire-tutorial-grpc/gen/todo/v1"
	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/config"
	appwire "github.com/97vaibhav/todo-crud-wire-tutorial-grpc/wire"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Wire resolves the entire dependency graph in one call.
	app, err := appwire.InitializeApp()
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	// Start the Kafka audit consumer in the background. It reads from todo-events
	// and writes to the audit_logs table. When you add Idea 2/3, start more consumers here.
	ctx := context.Background()
	go app.AuditConsumer.Run(ctx)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", cfg.GRPCPort, err)
	}

	// Register the auth interceptor so it runs before every handler.
	// grpc.ChainUnaryInterceptor lets you add more interceptors later (e.g. logging, metrics).
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			app.AuthInterceptor.Unary(),
		),
	)

	todov1.RegisterTodoServiceServer(grpcServer, app.TodoHandler)
	authv1.RegisterAuthServiceServer(grpcServer, app.AuthHandler)

	// Reflection lets grpcurl and Postman discover your API at runtime.
	reflection.Register(grpcServer)

	log.Printf("gRPC server listening on :%s", cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
