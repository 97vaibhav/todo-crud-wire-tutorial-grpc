package auth

import (
	"context"
	"errors"
	"time"

	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/config"
	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ── Claims ────────────────────────────────────────────────────────────────────

// UserClaims is the payload we embed in every JWT.
// It travels from the interceptor → handler via the request context.
type UserClaims struct {
	UserID    string
	Email     string
	GroupType domain.GroupType
}

// jwtClaims is the internal struct that maps to the JWT payload fields.
type jwtClaims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	GroupType string `json:"group_type"`
	jwt.RegisteredClaims
}

// ── Service interface ─────────────────────────────────────────────────────────

type JWTService interface {
	GenerateToken(userID, email string, groupType domain.GroupType) (string, error)
	ValidateToken(tokenStr string) (*UserClaims, error)
}

type jwtService struct {
	secret string
	expiry time.Duration
}

// NewJWTService is a Wire provider.
func NewJWTService(cfg *config.Config) JWTService {
	return &jwtService{
		secret: cfg.JWTSecret,
		expiry: 24 * time.Hour,
	}
}

func (s *jwtService) GenerateToken(userID, email string, groupType domain.GroupType) (string, error) {
	claims := &jwtClaims{
		UserID:    userID,
		Email:     email,
		GroupType: string(groupType),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func (s *jwtService) ValidateToken(tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwtClaims{}, func(t *jwt.Token) (any, error) {
		// Enforce that the signing method is exactly HS256 — reject any other algorithm.
		// Skipping this check is a well-known JWT vulnerability (algorithm confusion).
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return &UserClaims{
		UserID:    claims.UserID,
		Email:     claims.Email,
		GroupType: domain.GroupType(claims.GroupType),
	}, nil
}

// ── Context helpers ───────────────────────────────────────────────────────────

// contextKey is a private type so our key never collides with keys from other packages.
type contextKey string

const claimsKey contextKey = "user_claims"

// ContextWithClaims stores the validated claims in the request context.
// Called by the interceptor after a successful token validation.
func ContextWithClaims(ctx context.Context, claims *UserClaims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

// ClaimsFromContext retrieves claims that the interceptor stored.
// Returns a gRPC Unauthenticated error if the context has no claims.
func ClaimsFromContext(ctx context.Context) (*UserClaims, error) {
	claims, ok := ctx.Value(claimsKey).(*UserClaims)
	if !ok || claims == nil {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}
	return claims, nil
}
