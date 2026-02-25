package grpc

import (
	"context"

	todov1 "github.com/97vaibhav/todo-crud-wire-tutorial-grpc/gen/todo/v1"
	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/auth"
	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/domain"
	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TodoHandler struct {
	todov1.UnimplementedTodoServiceServer
	uc usecase.TodoUsecase
}

func NewTodoHandler(uc usecase.TodoUsecase) *TodoHandler {
	return &TodoHandler{uc: uc}
}

func (h *TodoHandler) CreateTodo(ctx context.Context, req *todov1.CreateTodoRequest) (*todov1.CreateTodoResponse, error) {
	// The interceptor already validated the JWT and stored claims in the context.
	// We never trust user_id from the request body — always derive it from the token.
	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	todo, err := h.uc.CreateTodo(claims.UserID, req.GetTitle(), req.GetDescription())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create todo: %v", err)
	}

	return &todov1.CreateTodoResponse{
		Todo: domainToProto(todo),
	}, nil
}

func (h *TodoHandler) ListTodos(ctx context.Context, in *todov1.ListTodosRequest) (*todov1.ListTodosResponse, error) {
	todos, err := h.uc.ListTodos()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list todos: %v", err)
	}
	return &todov1.ListTodosResponse{
		Todos: domainTodosToProtos(todos),
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

func domainTodosToProtos(todos []*domain.Todo) []*todov1.Todo {
	protoTodos := make([]*todov1.Todo, len(todos))
	for i, t := range todos {
		protoTodos[i] = domainToProto(t)
	}
	return protoTodos
}

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
