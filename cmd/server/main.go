package main

import (
	"fmt"
	"log"
	"net"

	todov1 "github.com/97vaibhav/todo-crud-wire-tutorial-grpc/gen/todo/v1"
	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/config"
	appwire "github.com/97vaibhav/todo-crud-wire-tutorial-grpc/wire"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load config early so we know the port before Wire runs.
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Wire builds the entire dependency graph: Config → DB → Repo → Usecase → Handler.
	// If any provider returns an error (e.g. DB connection fails), we exit here.
	todoHandler, err := appwire.InitializeTodoHandler()
	if err != nil {
		log.Fatalf("failed to initialize dependencies: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", cfg.GRPCPort, err)
	}

	grpcServer := grpc.NewServer()

	// Register our handler with the gRPC server.
	// This links the generated service interface to our concrete implementation.
	todov1.RegisterTodoServiceServer(grpcServer, todoHandler)

	// reflection lets grpcurl/Postman discover your API without a proto file.
	// Remove this in production if you don't want to expose the schema.
	reflection.Register(grpcServer)

	log.Printf("gRPC server listening on :%s", cfg.GRPCPort)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
