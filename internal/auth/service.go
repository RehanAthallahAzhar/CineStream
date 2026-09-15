package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"CineStream/internal/config"
	db "CineStream/internal/database/sqlc"
)

var (
	ErrEmailAlreadyExists = errors.New("email is already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type Service interface {
	Register(ctx context.Context, req RegisterRequest) (UserResponse, error)
	Login(ctx context.Context, req LoginRequest) (LoginResponse, error)
}

type service struct {
	repo Repository
	cfg  *config.Config
}

func NewService(repo Repository, cfg *config.Config) Service {
	return &service{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *service) Register(ctx context.Context, req RegisterRequest) (UserResponse, error) {
	existingUser, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err == nil && existingUser.ID != uuid.Nil {
		return UserResponse{}, ErrEmailAlreadyExists
	} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return UserResponse{}, fmt.Errorf("failed to check existing user: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return UserResponse{}, fmt.Errorf("failed to hash password: %w", err)
	}

	role := req.Role
	if role == "" {
		role = RoleCustomer
	}

	userID := uuid.New()

	params := db.CreateUserParams{
		ID:           userID,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FullName:     req.FullName,
		Role:         role,
	}

	createdUser, err := s.repo.CreateUser(ctx, params)
	if err != nil {
		return UserResponse{}, fmt.Errorf("failed to create user: %w", err)
	}

	return UserResponse{
		ID:        createdUser.ID.String(),
		Email:     createdUser.Email,
		FullName:  createdUser.FullName,
		Role:      createdUser.Role,
		CreatedAt: createdUser.CreatedAt.Time,
	}, nil
}

func (s *service) Login(ctx context.Context, req LoginRequest) (LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return LoginResponse{}, ErrInvalidCredentials
		}
		return LoginResponse{}, fmt.Errorf("failed to get user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return LoginResponse{}, ErrInvalidCredentials
	}
	expirationTime := time.Now().Add(time.Duration(s.cfg.JWTExpirationHours) * time.Hour)
	claims := jwt.MapClaims{
		"user_id":   user.ID.String(),
		"email":     user.Email,
		"role":      user.Role,
		"full_name": user.FullName,
		"exp":       expirationTime.Unix(),
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return LoginResponse{}, fmt.Errorf("failed to sign jwt token: %w", err)
	}

	return LoginResponse{
		Token: tokenString,
		User: UserResponse{
			ID:        user.ID.String(),
			Email:     user.Email,
			FullName:  user.FullName,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.Time,
		},
	}, nil
}
