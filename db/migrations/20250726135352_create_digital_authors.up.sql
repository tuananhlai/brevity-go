ALTER TABLE users DROP COLUMN IF EXISTS user_type;

CREATE TABLE IF NOT EXISTS digital_authors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    system_prompt VARCHAR(4000) NOT NULL,
    api_key_id UUID NOT NULL,
    avatar_url VARCHAR(1000),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY (owner_id) REFERENCES users(id),
    FOREIGN KEY (api_key_id) REFERENCES llm_api_keys(id)
);
COMMENT ON TABLE digital_authors IS 'A bot users which interacts with LLM providers to write new articles.';
COMMENT ON COLUMN digital_authors.owner_id IS 'The ID of the user who created this digital author.';
COMMENT ON COLUMN digital_authors.api_key_id IS 'The ID of an existing LLM provider API key.';

-- Remove the foreign key constraint on articles.author_id, since we intended to link the articles
-- with entities from `users` table
ALTER TABLE articles DROP CONSTRAINT IF EXISTS fk_articles_author;
COMMENT ON COLUMN articles.author_id IS 'The ID of the digital author / bot who wrote this article.';