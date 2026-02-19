package grpc

import (
	"context"

	authv1 "github.com/97vaibhav/todo-crud-wire-tutorial-grpc/gen/auth/v1"
	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/auth"
	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/domain"
	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthHandler struct {
	authv1.UnimplementedAuthServiceServer
	uc usecase.AuthUsecase
}

func NewAuthHandler(uc usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

// requireAdmin is a helper that reads claims from context and returns PermissionDenied
// if the caller is not an ADMIN. Called at the top of every admin-only method.
func requireAdmin(ctx context.Context) (*auth.UserClaims, error) {
	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if claims.GroupType != domain.GroupTypeAdmin {
		return nil, status.Error(codes.PermissionDenied, "admin access required")
	}
	return claims, nil
}

// ── Login (public) ────────────────────────────────────────────────────────────

func (h *AuthHandler) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	token, user, err := h.uc.Login(req.GetEmail(), req.GetPassword())
	if err != nil {
		// Return Unauthenticated — do NOT return Internal here.
		// Internal would leak that the email exists; Unauthenticated is intentionally vague.
		return nil, status.Errorf(codes.Unauthenticated, "%v", err)
	}

	return &authv1.LoginResponse{
		AccessToken: token,
		User:        domainUserToProto(user),
	}, nil
}

// ── CreateUser (admin only) ───────────────────────────────────────────────────

func (h *AuthHandler) CreateUser(ctx context.Context, req *authv1.CreateUserRequest) (*authv1.CreateUserResponse, error) {
	if _, err := requireAdmin(ctx); err != nil {
		return nil, err
	}

	user, err := h.uc.CreateUser(
		req.GetName(),
		req.GetEmail(),
		req.GetPassword(),
		req.GetGroupId(),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create user: %v", err)
	}

	return &authv1.CreateUserResponse{User: domainUserToProto(user)}, nil
}

// ── DeleteUser (admin only) ───────────────────────────────────────────────────

func (h *AuthHandler) DeleteUser(ctx context.Context, req *authv1.DeleteUserRequest) (*authv1.DeleteUserResponse, error) {
	if _, err := requireAdmin(ctx); err != nil {
		return nil, err
	}

	if err := h.uc.DeleteUser(req.GetUserId()); err != nil {
		return nil, status.Errorf(codes.Internal, "delete user: %v", err)
	}

	return &authv1.DeleteUserResponse{}, nil
}

// ── ListGroups (admin only) ───────────────────────────────────────────────────

func (h *AuthHandler) ListGroups(ctx context.Context, _ *authv1.ListGroupsRequest) (*authv1.ListGroupsResponse, error) {
	if _, err := requireAdmin(ctx); err != nil {
		return nil, err
	}

	groups, err := h.uc.ListGroups()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list groups: %v", err)
	}

	protoGroups := make([]*authv1.Group, len(groups))
	for i, g := range groups {
		protoGroups[i] = domainGroupToProto(g)
	}

	return &authv1.ListGroupsResponse{Groups: protoGroups}, nil
}

// ── Conversion helpers ────────────────────────────────────────────────────────

func domainUserToProto(u *domain.User) *authv1.User {
	return &authv1.User{
		Id:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		GroupId:   u.GroupID,
		GroupType: domainGroupTypeToProto(u.GroupType),
		CreatedAt: timestamppb.New(u.CreatedAt),
	}
}

func domainGroupToProto(g *domain.Group) *authv1.Group {
	return &authv1.Group{
		Id:        g.ID,
		Name:      g.Name,
		Type:      domainGroupTypeToProto(g.Type),
		CreatedAt: timestamppb.New(g.CreatedAt),
	}
}

func domainGroupTypeToProto(t domain.GroupType) authv1.GroupType {
	switch t {
	case domain.GroupTypeAdmin:
		return authv1.GroupType_GROUP_TYPE_ADMIN
	case domain.GroupTypeGuest:
		return authv1.GroupType_GROUP_TYPE_GUEST
	default:
		return authv1.GroupType_GROUP_TYPE_UNSPECIFIED
	}
}
