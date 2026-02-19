//go:build wireinject
// +build wireinject

package wire

import (
	grpchandler "github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/handler/grpc"
	infradb "github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/infrastructure/db"
	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/infrastructure/repository"
	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/auth"
	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/config"
	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/middleware"
	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/usecase"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	config.Load,
	auth.NewJWTService,
	infradb.NewPostgresDB,
	repository.NewTodoRepository,
	repository.NewUserRepository,
	repository.NewGroupRepository,
	usecase.NewTodoUsecase,
	usecase.NewAuthUsecase,
	grpchandler.NewTodoHandler,
	grpchandler.NewAuthHandler,
	middleware.NewAuthInterceptor,
	newApp,
)

// InitializeApp is the single injector function Wire reads.
// Wire traces all dependencies through ProviderSet and writes wire_gen.go.
func InitializeApp() (*App, error) {
	wire.Build(ProviderSet)
	return nil, nil
}
