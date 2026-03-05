package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	authv1 "github.com/97vaibhav/todo-crud-wire-tutorial-grpc/gen/auth/v1"
	todov1 "github.com/97vaibhav/todo-crud-wire-tutorial-grpc/gen/todo/v1"
	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/config"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Forward HTTP "Authorization" header to gRPC metadata "authorization" so the
	// backend's auth interceptor can validate the JWT.
	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
			if strings.EqualFold(key, "Authorization") {
				return "authorization", true
			}
			return runtime.DefaultHeaderMatcher(key)
		}),
	)

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := todov1.RegisterTodoServiceHandlerFromEndpoint(ctx, mux, cfg.GRPCBackendAddr, opts); err != nil {
		log.Fatalf("register todo service gateway: %v", err)
	}
	if err := authv1.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, cfg.GRPCBackendAddr, opts); err != nil {
		log.Fatalf("register auth service gateway: %v", err)
	}

	addr := ":" + cfg.GatewayPort
	log.Printf("REST API gateway listening on %s (backend gRPC: %s)", addr, cfg.GRPCBackendAddr)
	if err := http.ListenAndServe(addr, mux); err != nil && err != http.ErrServerClosed {
		log.Fatalf("gateway serve: %v", err)
	}
}
