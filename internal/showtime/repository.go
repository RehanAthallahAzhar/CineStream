package showtime

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"

	db "CineStream/internal/database/sqlc"
)

type Repository interface {
	CreateShowtime(ctx context.Context, arg db.CreateShowtimeParams) (db.Showtime, error)
	GetShowtimeByID(ctx context.Context, id uuid.UUID) (db.GetShowtimeByIDRow, error)
	ListShowtimes(ctx context.Context, limit, offset int32) ([]db.ListShowtimesRow, error)
	UpdateShowtime(ctx context.Context, arg db.UpdateShowtimeParams) (db.Showtime, error)
	DeleteShowtime(ctx context.Context, id uuid.UUID) error
	CheckShowtimeOverlap(ctx context.Context, studioID, id uuid.UUID, startTime, endTime string) (int64, error)
}

type repository struct {
	queries *db.Queries
}

func NewRepository(sqlDB *sql.DB) Repository {
	return &repository{
		queries: db.New(sqlDB),
	}
}

func (r *repository) CreateShowtime(ctx context.Context, arg db.CreateShowtimeParams) (db.Showtime, error) {
	return r.queries.CreateShowtime(ctx, arg)
}

func (r *repository) GetShowtimeByID(ctx context.Context, id uuid.UUID) (db.GetShowtimeByIDRow, error) {
	return r.queries.GetShowtimeByID(ctx, id)
}

func (r *repository) ListShowtimes(ctx context.Context, limit, offset int32) ([]db.ListShowtimesRow, error) {
	return r.queries.ListShowtimes(ctx, db.ListShowtimesParams{
		Limit:  limit,
		Offset: offset,
	})
}

func (r *repository) UpdateShowtime(ctx context.Context, arg db.UpdateShowtimeParams) (db.Showtime, error) {
	return r.queries.UpdateShowtime(ctx, arg)
}

func (r *repository) DeleteShowtime(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteShowtime(ctx, id)
}

func (r *repository) CheckShowtimeOverlap(ctx context.Context, studioID, id uuid.UUID, startTime, endTime string) (int64, error) {
	st, _ := time.Parse(time.RFC3339, startTime)
	et, _ := time.Parse(time.RFC3339, endTime)

	return r.queries.CheckShowtimeOverlap(ctx, db.CheckShowtimeOverlapParams{
		StudioID:  studioID,
		ID:        id,
		StartTime: st,
		EndTime:   et,
	})
}
