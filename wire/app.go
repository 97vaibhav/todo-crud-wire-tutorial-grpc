package wire

import (
	grpchandler "github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/handler/grpc"
	kafkaconsumer "github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/kafka/consumer"
	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/middleware"
)

// App bundles everything main.go needs after the dependency graph is wired.
type App struct {
	TodoHandler     *grpchandler.TodoHandler
	AuthHandler     *grpchandler.AuthHandler
	AuthInterceptor *middleware.AuthInterceptor
	AuditConsumer   *kafkaconsumer.AuditConsumer
}

func newApp(
	todo *grpchandler.TodoHandler,
	auth *grpchandler.AuthHandler,
	interceptor *middleware.AuthInterceptor,
	auditConsumer *kafkaconsumer.AuditConsumer,
) *App {
	return &App{
		TodoHandler:     todo,
		AuthHandler:     auth,
		AuthInterceptor: interceptor,
		AuditConsumer:   auditConsumer,
	}
}
