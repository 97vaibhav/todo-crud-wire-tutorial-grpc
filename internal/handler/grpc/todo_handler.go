package grpc

import (
	"context"

	todov1 "github.com/97vaibhav/todo-crud-wire-tutorial-grpc/gen/todo/v1"
	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/domain"
	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TodoHandler implements the generated todov1.TodoServiceServer interface.
// Embedding UnimplementedTodoServiceServer means any RPC you haven't implemented yet
// returns codes.Unimplemented instead of crashing — future-proof by design.
type TodoHandler struct {
	todov1.UnimplementedTodoServiceServer
	uc usecase.TodoUsecase
}

// NewTodoHandler is a Wire provider.
func NewTodoHandler(uc usecase.TodoUsecase) *TodoHandler {
	return &TodoHandler{uc: uc}
}

func (h *TodoHandler) CreateTodo(ctx context.Context, req *todov1.CreateTodoRequest) (*todov1.CreateTodoResponse, error) {
	todo, err := h.uc.CreateTodo(req.GetTitle(), req.GetDescription())
	if err != nil {
		// Wrap Go errors into proper gRPC status errors.
		// The client receives a structured error with a status code, not a raw string.
		return nil, status.Errorf(codes.Internal, "create todo: %v", err)
	}

	return &todov1.CreateTodoResponse{
		Todo: domainToProto(todo),
	}, nil
}

func (h *TodoHandler) GetTodo(ctx context.Context, req *todov1.GetTodoRequest) (*todov1.GetTodoResponse, error) {
	todo, err := h.uc.GetTodo(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get todo: %v", err)
	}
	return &todov1.GetTodoResponse{
		Todo: domainToProto(todo),
	}, nil
}

// domainToProto converts a domain.Todo → proto Todo message.
// Keeping this as a separate function makes it reusable when you add GetTodo, ListTodos, etc.
func domainToProto(t *domain.Todo) *todov1.Todo {
	return &todov1.Todo{
		Id:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Status:      domainStatusToProto(t.Status),
		CreatedAt:   timestamppb.New(t.CreatedAt),
		UpdatedAt:   timestamppb.New(t.UpdatedAt),
	}
}

func domainStatusToProto(s domain.TodoStatus) todov1.TodoStatus {
	switch s {
	case domain.TodoStatusPending:
		return todov1.TodoStatus_TODO_STATUS_PENDING
	case domain.TodoStatusInProgress:
		return todov1.TodoStatus_TODO_STATUS_IN_PROGRESS
	case domain.TodoStatusDone:
		return todov1.TodoStatus_TODO_STATUS_DONE
	default:
		return todov1.TodoStatus_TODO_STATUS_UNSPECIFIED
	}
}
