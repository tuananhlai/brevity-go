COMMENT ON COLUMN articles.author_id IS NULL;
ALTER TABLE articles ADD CONSTRAINT fk_articles_author FOREIGN KEY (author_id) REFERENCES users(id);
ALTER TABLE users ADD COLUMN IF NOT EXISTS user_type VARCHAR(15) NOT NULL DEFAULT 'user';
DROP TABLE IF EXISTS digital_authors;
