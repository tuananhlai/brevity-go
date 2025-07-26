ALTER TABLE users DROP COLUMN IF EXISTS user_type;

CREATE TABLE IF NOT EXISTS digital_authors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    system_prompt VARCHAR(4000) NOT NULL,
    avatar_url VARCHAR(1000),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY (owner_id) REFERENCES users(id)
);
COMMENT ON COLUMN digital_authors.owner_id IS 'The ID of the user who created this digital author.';