-- +migrate Up
CREATE TABLE IF NOT EXISTS digital_authors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    display_name VARCHAR(255) NOT NULL,    
)