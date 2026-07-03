CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TABLE medicines (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_medicines_name_trgm ON medicines USING GIN (name gin_trgm_ops);
CREATE INDEX idx_medicines_name ON medicines (name);
