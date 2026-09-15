BEGIN;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- ENUM TYPES 
CREATE TYPE seat_type_enum        AS ENUM ('regular', 'vip', 'couple');
CREATE TYPE studio_layout_enum    AS ENUM ('2D', '3D', 'IMAX', '4DX');
CREATE TYPE showtime_status_enum  AS ENUM ('SCHEDULED', 'ONGOING', 'CANCELLED', 'DONE');
CREATE TYPE showtime_seat_status_enum AS ENUM ('AVAILABLE', 'LOCKED', 'SOLD', 'VOID');
CREATE TYPE booking_status_enum   AS ENUM ('PENDING', 'PAID', 'EXPIRED', 'CANCELLED', 'REFUND_PENDING', 'REFUNDED');
CREATE TYPE payment_method_enum   AS ENUM ('va', 'ewallet', 'card', 'qris');
CREATE TYPE payment_status_enum   AS ENUM ('PENDING', 'SUCCESS', 'FAILED', 'EXPIRED');
CREATE TYPE ticket_status_enum    AS ENUM ('ACTIVE', 'USED', 'VOID', 'REFUNDED');
CREATE TYPE refund_reason_enum    AS ENUM ('customer_cancel', 'cinema_cancel', 'reschedule');
CREATE TYPE refund_status_enum    AS ENUM ('PENDING', 'PROCESSING', 'SUCCESS', 'FAILED');

CREATE TABLE city (
    city_id     BIGSERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE cinema (
    cinema_id   BIGSERIAL PRIMARY KEY,
    city_id     BIGINT NOT NULL REFERENCES city(city_id) ON DELETE RESTRICT,
    name        VARCHAR(150) NOT NULL,
    address     VARCHAR(255) NOT NULL,
    phone       VARCHAR(30),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_cinema_city_id ON cinema(city_id);

CREATE TABLE studio (
    studio_id   BIGSERIAL PRIMARY KEY,
    cinema_id   BIGINT NOT NULL REFERENCES cinema(cinema_id) ON DELETE RESTRICT,
    name        VARCHAR(50) NOT NULL,
    total_seat  INT NOT NULL CHECK (total_seat > 0),
    layout_type studio_layout_enum NOT NULL DEFAULT '2D',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (cinema_id, name)
);
CREATE INDEX idx_studio_cinema_id ON studio(cinema_id);

CREATE TABLE seat (
    seat_id     BIGSERIAL PRIMARY KEY,
    studio_id   BIGINT NOT NULL REFERENCES studio(studio_id) ON DELETE CASCADE,
    seat_row    VARCHAR(5) NOT NULL,
    seat_number VARCHAR(5) NOT NULL,
    seat_type   seat_type_enum NOT NULL DEFAULT 'regular',
    is_active   BOOLEAN NOT NULL DEFAULT TRUE,
    UNIQUE (studio_id, seat_row, seat_number)
);
CREATE INDEX idx_seat_studio_id ON seat(studio_id);

CREATE TABLE movie (
    movie_id         BIGSERIAL PRIMARY KEY,
    title            VARCHAR(200) NOT NULL,
    duration_minutes INT NOT NULL CHECK (duration_minutes > 0),
    rating           VARCHAR(10) NOT NULL, -- SU/13+/17+/21+
    genre            VARCHAR(100),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE showtime (
    showtime_id BIGSERIAL PRIMARY KEY,
    movie_id    BIGINT NOT NULL REFERENCES movie(movie_id) ON DELETE RESTRICT,
    studio_id   BIGINT NOT NULL REFERENCES studio(studio_id) ON DELETE RESTRICT,
    start_time  TIMESTAMPTZ NOT NULL,
    end_time    TIMESTAMPTZ NOT NULL,
    base_price  DECIMAL(12,2) NOT NULL CHECK (base_price >= 0),
    status      showtime_status_enum NOT NULL DEFAULT 'SCHEDULED',
    version     INT NOT NULL DEFAULT 0,  -- optimistic locking
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (end_time > start_time)
);
CREATE INDEX idx_showtime_studio_time ON showtime(studio_id, start_time);
CREATE INDEX idx_showtime_movie_id ON showtime(movie_id);

CREATE TABLE showtime_seat (
    showtime_seat_id  BIGSERIAL PRIMARY KEY,
    showtime_id        BIGINT NOT NULL REFERENCES showtime(showtime_id) ON DELETE CASCADE,
    seat_id             BIGINT NOT NULL REFERENCES seat(seat_id) ON DELETE RESTRICT,
    status              showtime_seat_status_enum NOT NULL DEFAULT 'AVAILABLE',
    lock_expires_at     TIMESTAMPTZ,
    locked_by_session   VARCHAR(100),
    version             INT NOT NULL DEFAULT 0,
    UNIQUE (showtime_id, seat_id)
);
-- Index
CREATE INDEX idx_showtime_seat_showtime_status ON showtime_seat(showtime_id, status);
CREATE INDEX idx_showtime_seat_lock_expiry ON showtime_seat(lock_expires_at)
    WHERE status = 'LOCKED'; -- partial index

CREATE TABLE app_user ( 
    user_id        BIGSERIAL PRIMARY KEY,
    name            VARCHAR(150) NOT NULL,
    email           VARCHAR(150) NOT NULL UNIQUE,
    phone           VARCHAR(30),
    password_hash   VARCHAR(255) NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE booking (
    booking_id      BIGSERIAL PRIMARY KEY,
    user_id          BIGINT NOT NULL REFERENCES app_user(user_id) ON DELETE RESTRICT,
    showtime_id      BIGINT NOT NULL REFERENCES showtime(showtime_id) ON DELETE RESTRICT,
    booking_code     VARCHAR(20) NOT NULL UNIQUE,
    status           booking_status_enum NOT NULL DEFAULT 'PENDING',
    total_amount     DECIMAL(12,2) NOT NULL CHECK (total_amount >= 0),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at       TIMESTAMPTZ NOT NULL
);
CREATE INDEX idx_booking_user_id ON booking(user_id);
CREATE INDEX idx_booking_showtime_id ON booking(showtime_id);
CREATE INDEX idx_booking_status_expiry ON booking(status, expires_at)
    WHERE status = 'PENDING';

CREATE TABLE booking_seat (
    booking_seat_id   BIGSERIAL PRIMARY KEY,
    booking_id         BIGINT NOT NULL REFERENCES booking(booking_id) ON DELETE CASCADE,
    showtime_seat_id   BIGINT NOT NULL REFERENCES showtime_seat(showtime_seat_id) ON DELETE RESTRICT,
    price              DECIMAL(12,2) NOT NULL CHECK (price >= 0),
    UNIQUE (showtime_seat_id)
);
CREATE INDEX idx_booking_seat_booking_id ON booking_seat(booking_id);

CREATE TABLE payment (
    payment_id        BIGSERIAL PRIMARY KEY,
    booking_id          BIGINT NOT NULL UNIQUE REFERENCES booking(booking_id) ON DELETE RESTRICT,
    method              payment_method_enum NOT NULL,
    status              payment_status_enum NOT NULL DEFAULT 'PENDING',
    amount              DECIMAL(12,2) NOT NULL CHECK (amount >= 0),
    paid_at             TIMESTAMPTZ,
    gateway_ref         VARCHAR(100) UNIQUE,
    idempotency_key     VARCHAR(100) NOT NULL UNIQUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_payment_status ON payment(status);

CREATE TABLE ticket (
    ticket_id          BIGSERIAL PRIMARY KEY,
    booking_seat_id      BIGINT NOT NULL UNIQUE REFERENCES booking_seat(booking_seat_id) ON DELETE RESTRICT,
    qr_code              VARCHAR(150) NOT NULL UNIQUE,
    status               ticket_status_enum NOT NULL DEFAULT 'ACTIVE',
    issued_at            TIMESTAMPTZ NOT NULL DEFAULT now(),
    used_at              TIMESTAMPTZ
);

CREATE TABLE refund (
    refund_id       BIGSERIAL PRIMARY KEY,
    booking_id        BIGINT NOT NULL REFERENCES booking(booking_id) ON DELETE RESTRICT,
    payment_id        BIGINT NOT NULL REFERENCES payment(payment_id) ON DELETE RESTRICT,
    reason            refund_reason_enum NOT NULL,
    status            refund_status_enum NOT NULL DEFAULT 'PENDING',
    amount            DECIMAL(12,2) NOT NULL CHECK (amount >= 0),
    refund_ref        VARCHAR(100),
    requested_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    processed_at      TIMESTAMPTZ
);
CREATE INDEX idx_refund_booking_id ON refund(booking_id);
CREATE INDEX idx_refund_status ON refund(status);

COMMIT;
