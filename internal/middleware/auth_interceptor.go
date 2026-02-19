package middleware

import (
	"context"
	"strings"

	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// publicMethods lists gRPC full method names that do NOT require a JWT.
// Format: /<proto_package>.<Service>/<Method>
var publicMethods = map[string]bool{
	"/auth.v1.AuthService/Login": true,
}

// AuthInterceptor is a struct wrapper around the gRPC unary interceptor function.
// Wrapping it lets Wire treat it as a named dependency rather than a raw function type.
type AuthInterceptor struct {
	jwtSvc auth.JWTService
}

// NewAuthInterceptor is a Wire provider.
func NewAuthInterceptor(jwtSvc auth.JWTService) *AuthInterceptor {
	return &AuthInterceptor{jwtSvc: jwtSvc}
}

// Unary returns the grpc.UnaryServerInterceptor to be registered with the gRPC server.
// It is called once at startup, not per-request.
func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		// Skip auth entirely for public endpoints.
		if publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		// gRPC metadata is the equivalent of HTTP headers.
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		// Convention: "authorization" header with value "Bearer <token>"
		values := md["authorization"]
		if len(values) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization header")
		}

		tokenStr := strings.TrimPrefix(values[0], "Bearer ")
		if tokenStr == values[0] {
			// TrimPrefix didn't change anything — "Bearer " prefix was missing.
			return nil, status.Error(codes.Unauthenticated, "authorization header must start with 'Bearer '")
		}

		claims, err := i.jwtSvc.ValidateToken(tokenStr)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		// Store claims in context so handlers can read them with auth.ClaimsFromContext.
		ctx = auth.ContextWithClaims(ctx, claims)
		return handler(ctx, req)
	}
}
