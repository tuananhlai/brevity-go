-- +migrate Up
BEGIN;

ALTER TABLE users DROP COLUMN IF EXISTS user_type;

CREATE TABLE IF NOT EXISTS digital_authors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    display_name VARCHAR(255) NOT NULL,    
    system_prompt VARCHAR(4000) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

COMMENT ON TABLE digital_authors IS 'A bot users which interacts with LLM providers to write new articles.';

-- Remove the foreign key constraint on articles.author_id, since we intended to link the articles
-- with entities from `users` table
ALTER TABLE articles DROP CONSTRAINT IF EXISTS fk_articles_author;
ALTER TABLE articles ADD CONSTRAINT fk_articles_digital_authors FOREIGN KEY (author_id) REFERENCES digital_authors(id);
COMMENT ON COLUMN articles.author_id IS 'The ID of the digital author / bot who wrote this article.';

COMMIT;