package auth

import (
	"context"
	"database/sql"

	db "CineStream/internal/database/sqlc"

	"github.com/google/uuid"
)

type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (db.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (db.User, error)
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.CreateUserRow, error)
}

type repository struct {
	queries *db.Queries
}

func NewRepository(sqlDB *sql.DB) Repository {
	return &repository{
		queries: db.New(sqlDB),
	}
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	return r.queries.GetUserByEmail(ctx, email)
}

func (r *repository) GetUserByID(ctx context.Context, id uuid.UUID) (db.User, error) {
	return r.queries.GetUserByID(ctx, id)
}

func (r *repository) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.CreateUserRow, error) {
	return r.queries.CreateUser(ctx, arg)
}
