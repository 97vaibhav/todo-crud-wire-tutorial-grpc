//go:build wireinject
// +build wireinject

// The build tag above means this file is ONLY read by the wire CLI tool,
// not by the normal Go compiler. wire_gen.go is read by the Go compiler instead.

package wire

import (
	grpchandler "github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/handler/grpc"
	infradb "github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/infrastructure/db"
	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/infrastructure/repository"
	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/config"
	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/usecase"
	"github.com/google/wire"
)

// ProviderSet groups all providers. Wire resolves the dependency graph from this set.
// Order does NOT matter here — Wire figures out the order automatically.
var ProviderSet = wire.NewSet(
	config.Load,
	infradb.NewPostgresDB,
	repository.NewTodoRepository,
	usecase.NewTodoUsecase,
	grpchandler.NewTodoHandler,
)

// InitializeTodoHandler is the "injector" function.
// Wire reads its signature: "I want a *TodoHandler and I have no inputs".
// Wire then traces the dependency graph through ProviderSet and generates
// the wiring code in wire_gen.go automatically.
func InitializeTodoHandler() (*grpchandler.TodoHandler, error) {
	wire.Build(ProviderSet)
	return nil, nil // Wire replaces this body entirely in wire_gen.go
}
