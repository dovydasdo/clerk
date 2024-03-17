CREATE TABLE locations (
    id INTEGER PRIMARY KEY,
    created_at TIMESTAMP NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,
    lat NUMERIC NULL,
    lng NUMERIC NULL
);

CREATE INDEX IF NOT EXISTS idx_locations_deleted_at ON locations (deleted_at);

