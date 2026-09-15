package showtime

import "time"

type CreateShowtimeRequest struct {
	MovieID   string    `json:"movie_id" validate:"required,uuid"`
	StudioID  string    `json:"studio_id" validate:"required,uuid"`
	StartTime time.Time `json:"start_time" validate:"required"`
	EndTime   time.Time `json:"end_time" validate:"required"`
	Price     float64   `json:"price" validate:"required,gt=0"`
	Status    string    `json:"status" validate:"omitempty,oneof=SCHEDULED CANCELLED COMPLETED"`
}

type UpdateShowtimeRequest struct {
	MovieID   string    `json:"movie_id" validate:"required,uuid"`
	StudioID  string    `json:"studio_id" validate:"required,uuid"`
	StartTime time.Time `json:"start_time" validate:"required"`
	EndTime   time.Time `json:"end_time" validate:"required"`
	Price     float64   `json:"price" validate:"required,gt=0"`
	Status    string    `json:"status" validate:"required,oneof=SCHEDULED CANCELLED COMPLETED"`
}

type ShowtimeResponse struct {
	ID            string    `json:"id"`
	MovieID       string    `json:"movie_id"`
	StudioID      string    `json:"studio_id"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	Price         string    `json:"price"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	MovieTitle    string    `json:"movie_title,omitempty"`
	MovieDuration int32     `json:"movie_duration,omitempty"`
	MovieRating   string    `json:"movie_rating,omitempty"`
	StudioName    string    `json:"studio_name,omitempty"`
	CinemaName    string    `json:"cinema_name,omitempty"`
}

type ListShowtimesQuery struct {
	Page  int `query:"page"`
	Limit int `query:"limit"`
}

type ListShowtimesResponse struct {
	Items []ShowtimeResponse `json:"items"`
	Page  int                `json:"page"`
	Limit int                `json:"limit"`
}
