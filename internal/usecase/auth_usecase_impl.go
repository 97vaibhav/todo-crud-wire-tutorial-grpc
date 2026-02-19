package usecase

import (
	"errors"
	"strings"

	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/auth"
	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	userRepo  domain.UserRepository
	groupRepo domain.GroupRepository
	jwtSvc    auth.JWTService
}

// NewAuthUsecase is a Wire provider.
func NewAuthUsecase(
	userRepo domain.UserRepository,
	groupRepo domain.GroupRepository,
	jwtSvc auth.JWTService,
) AuthUsecase {
	return &authUsecase{
		userRepo:  userRepo,
		groupRepo: groupRepo,
		jwtSvc:    jwtSvc,
	}
}

func (u *authUsecase) Login(email, password string) (string, *domain.User, error) {
	user, err := u.userRepo.FindByEmail(strings.ToLower(email))
	if err != nil {
		// Don't reveal whether the email exists — return a generic message.
		return "", nil, errors.New("invalid email or password")
	}

	// bcrypt.CompareHashAndPassword handles the salt automatically.
	// It returns an error if the password doesn't match.
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	token, err := u.jwtSvc.GenerateToken(user.ID, user.Email, user.GroupType)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (u *authUsecase) CreateUser(name, email, password, groupID string) (*domain.User, error) {
	if name == "" || email == "" || password == "" || groupID == "" {
		return nil, errors.New("name, email, password, and group_id are all required")
	}

	// Validate the group exists before creating the user.
	group, err := u.groupRepo.FindByID(groupID)
	if err != nil {
		return nil, errors.New("group not found")
	}

	// Cost 12 is a good balance: ~300ms on modern hardware, painful to brute-force.
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Name:         name,
		Email:        strings.ToLower(email),
		PasswordHash: string(hash),
		GroupID:      groupID,
		GroupType:    group.Type,
	}

	return u.userRepo.Create(user)
}

func (u *authUsecase) DeleteUser(userID string) error {
	if _, err := u.userRepo.FindByID(userID); err != nil {
		return errors.New("user not found")
	}
	return u.userRepo.Delete(userID)
}

func (u *authUsecase) ListGroups() ([]*domain.Group, error) {
	return u.groupRepo.List()
}
