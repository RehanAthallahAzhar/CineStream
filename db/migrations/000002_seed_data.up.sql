-- Seed Initial Users (Password for both: password123)
-- bcrypt hash for 'password123': $2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy
INSERT INTO users (id, email, password_hash, full_name, role)
VALUES 
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'admin@cinestream.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Super Admin', 'ADMIN'),
    ('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'customer@cinestream.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'John Doe', 'CUSTOMER')
ON CONFLICT (email) DO NOTHING;

-- Seed Cinema
INSERT INTO cinemas (id, name, city, address)
VALUES 
    ('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'Grand Cinema XXI', 'Jakarta', 'Jl. M.H. Thamrin No. 1')
ON CONFLICT (id) DO NOTHING;

-- Seed Studio
INSERT INTO studios (id, cinema_id, name, total_seats)
VALUES 
    ('d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'Studio 1 IMAX', 100)
ON CONFLICT (id) DO NOTHING;

-- Seed Movie
INSERT INTO movies (id, title, duration_minutes, rating)
VALUES 
    ('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a55', 'Inception 2', 148, '13+')
ON CONFLICT (id) DO NOTHING;