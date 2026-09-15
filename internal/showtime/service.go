package showtime

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"

	db "CineStream/internal/database/sqlc"
)

var (
	ErrShowtimeNotFound  = errors.New("showtime not found")
	ErrInvalidTimeRange  = errors.New("end_time must be after start_time")
	ErrScheduleConflict  = errors.New("schedule overlaps with an existing showtime in the same studio")
	ErrInvalidUUIDFormat = errors.New("invalid uuid format")
)

type Service interface {
	Create(ctx context.Context, req CreateShowtimeRequest) (ShowtimeResponse, error)
	GetByID(ctx context.Context, id string) (ShowtimeResponse, error)
	List(ctx context.Context, query ListShowtimesQuery) (ListShowtimesResponse, error)
	Update(ctx context.Context, id string, req UpdateShowtimeRequest) (ShowtimeResponse, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(ctx context.Context, req CreateShowtimeRequest) (ShowtimeResponse, error) {
	if !req.EndTime.After(req.StartTime) {
		return ShowtimeResponse{}, ErrInvalidTimeRange
	}

	movieUUID, err := uuid.Parse(req.MovieID)
	if err != nil {
		return ShowtimeResponse{}, fmt.Errorf("%w: invalid movie_id", ErrInvalidUUIDFormat)
	}

	studioUUID, err := uuid.Parse(req.StudioID)
	if err != nil {
		return ShowtimeResponse{}, fmt.Errorf("%w: invalid studio_id", ErrInvalidUUIDFormat)
	}

	newShowtimeID := uuid.New()

	// Check schedule overlap
	overlapCount, err := s.repo.CheckShowtimeOverlap(
		ctx,
		studioUUID,
		newShowtimeID,
		req.StartTime.Format(time.RFC3339),
		req.EndTime.Format(time.RFC3339),
	)
	if err != nil {
		return ShowtimeResponse{}, fmt.Errorf("failed to check schedule overlap: %w", err)
	}
	if overlapCount > 0 {
		return ShowtimeResponse{}, ErrScheduleConflict
	}

	status := req.Status
	if status == "" {
		status = StatusScheduled
	}

	priceStr := strconv.FormatFloat(req.Price, 'f', 2, 64)

	params := db.CreateShowtimeParams{
		ID:        newShowtimeID,
		MovieID:   movieUUID,
		StudioID:  studioUUID,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Price:     priceStr,
		Status:    status,
	}

	res, err := s.repo.CreateShowtime(ctx, params)
	if err != nil {
		return ShowtimeResponse{}, fmt.Errorf("failed to create showtime: %w", err)
	}

	return ShowtimeResponse{
		ID:        res.ID.String(),
		MovieID:   res.MovieID.String(),
		StudioID:  res.StudioID.String(),
		StartTime: res.StartTime,
		EndTime:   res.EndTime,
		Price:     res.Price,
		Status:    res.Status,
		CreatedAt: res.CreatedAt.Time,
		UpdatedAt: res.UpdatedAt.Time,
	}, nil
}

func (s *service) GetByID(ctx context.Context, id string) (ShowtimeResponse, error) {
	showtimeUUID, err := uuid.Parse(id)
	if err != nil {
		return ShowtimeResponse{}, ErrInvalidUUIDFormat
	}

	res, err := s.repo.GetShowtimeByID(ctx, showtimeUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ShowtimeResponse{}, ErrShowtimeNotFound
		}
		return ShowtimeResponse{}, fmt.Errorf("failed to get showtime: %w", err)
	}

	return ShowtimeResponse{
		ID:            res.ID.String(),
		MovieID:       res.MovieID.String(),
		StudioID:      res.StudioID.String(),
		StartTime:     res.StartTime,
		EndTime:       res.EndTime,
		Price:         res.Price,
		Status:        res.Status,
		CreatedAt:     res.CreatedAt.Time,
		UpdatedAt:     res.UpdatedAt.Time,
		MovieTitle:    res.MovieTitle,
		MovieDuration: res.MovieDuration,
		MovieRating:   res.MovieRating,
		StudioName:    res.StudioName,
		CinemaName:    res.CinemaName,
	}, nil
}

func (s *service) List(ctx context.Context, query ListShowtimesQuery) (ListShowtimesResponse, error) {
	page := query.Page
	if page < 1 {
		page = 1
	}

	limit := query.Limit
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := int32((page - 1) * limit)

	rows, err := s.repo.ListShowtimes(ctx, int32(limit), offset)
	if err != nil {
		return ListShowtimesResponse{}, fmt.Errorf("failed to list showtimes: %w", err)
	}

	items := make([]ShowtimeResponse, 0, len(rows))
	for _, r := range rows {
		items = append(items, ShowtimeResponse{
			ID:            r.ID.String(),
			MovieID:       r.MovieID.String(),
			StudioID:      r.StudioID.String(),
			StartTime:     r.StartTime,
			EndTime:       r.EndTime,
			Price:         r.Price,
			Status:        r.Status,
			CreatedAt:     r.CreatedAt.Time,
			UpdatedAt:     r.UpdatedAt.Time,
			MovieTitle:    r.MovieTitle,
			MovieDuration: r.MovieDuration,
			MovieRating:   r.MovieRating,
			StudioName:    r.StudioName,
			CinemaName:    r.CinemaName,
		})
	}

	return ListShowtimesResponse{
		Items: items,
		Page:  page,
		Limit: limit,
	}, nil
}

func (s *service) Update(ctx context.Context, id string, req UpdateShowtimeRequest) (ShowtimeResponse, error) {
	showtimeUUID, err := uuid.Parse(id)
	if err != nil {
		return ShowtimeResponse{}, ErrInvalidUUIDFormat
	}

	if !req.EndTime.After(req.StartTime) {
		return ShowtimeResponse{}, ErrInvalidTimeRange
	}

	movieUUID, err := uuid.Parse(req.MovieID)
	if err != nil {
		return ShowtimeResponse{}, fmt.Errorf("%w: invalid movie_id", ErrInvalidUUIDFormat)
	}

	studioUUID, err := uuid.Parse(req.StudioID)
	if err != nil {
		return ShowtimeResponse{}, fmt.Errorf("%w: invalid studio_id", ErrInvalidUUIDFormat)
	}

	// Check if showtime exists
	_, err = s.repo.GetShowtimeByID(ctx, showtimeUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ShowtimeResponse{}, ErrShowtimeNotFound
		}
		return ShowtimeResponse{}, err
	}

	// Check schedule overlap excluding self
	overlapCount, err := s.repo.CheckShowtimeOverlap(
		ctx,
		studioUUID,
		showtimeUUID,
		req.StartTime.Format(time.RFC3339),
		req.EndTime.Format(time.RFC3339),
	)
	if err != nil {
		return ShowtimeResponse{}, fmt.Errorf("failed to check schedule overlap: %w", err)
	}
	if overlapCount > 0 {
		return ShowtimeResponse{}, ErrScheduleConflict
	}

	priceStr := strconv.FormatFloat(req.Price, 'f', 2, 64)

	params := db.UpdateShowtimeParams{
		ID:        showtimeUUID,
		MovieID:   movieUUID,
		StudioID:  studioUUID,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Price:     priceStr,
		Status:    req.Status,
	}

	res, err := s.repo.UpdateShowtime(ctx, params)
	if err != nil {
		return ShowtimeResponse{}, fmt.Errorf("failed to update showtime: %w", err)
	}

	return ShowtimeResponse{
		ID:        res.ID.String(),
		MovieID:   res.MovieID.String(),
		StudioID:  res.StudioID.String(),
		StartTime: res.StartTime,
		EndTime:   res.EndTime,
		Price:     res.Price,
		Status:    res.Status,
		CreatedAt: res.CreatedAt.Time,
		UpdatedAt: res.UpdatedAt.Time,
	}, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	showtimeUUID, err := uuid.Parse(id)
	if err != nil {
		return ErrInvalidUUIDFormat
	}

	_, err = s.repo.GetShowtimeByID(ctx, showtimeUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrShowtimeNotFound
		}
		return err
	}

	if err := s.repo.DeleteShowtime(ctx, showtimeUUID); err != nil {
		return fmt.Errorf("failed to delete showtime: %w", err)
	}

	return nil
}
