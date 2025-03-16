CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    username VARCHAR(15) NOT NULL,
    user_type VARCHAR(15) NOT NULL DEFAULT 'user',
    email VARCHAR(255),
    password_hash VARCHAR(255),
    display_name VARCHAR(255),
    avatar_url VARCHAR(1000),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE (username),
    UNIQUE (email)
);
COMMENT ON COLUMN users.email IS 'The email address of the user. It is required for human users and optional for bots.';
COMMENT ON COLUMN users.user_type IS 'The type of user, either "user" or "bot".';
COMMENT ON COLUMN users.display_name IS 'The user name to display in the UI. If not set, the username will be used.';
COMMENT ON COLUMN users.password_hash IS 'A hashed password string using bcrypt.';
CREATE TABLE IF NOT EXISTS articles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    slug VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(500),
    text_content TEXT NOT NULL,
    content TEXT NOT NULL,
    author_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_articles_author FOREIGN KEY (author_id) REFERENCES users(id),
    UNIQUE (slug)
);
COMMENT ON COLUMN articles.slug IS 'The slug of the article, used for the URL. Example: "my-article-slug-3921"';
COMMENT ON COLUMN articles.description IS 'A short description of the article, provided by the author, used for previewing the article content.';
COMMENT ON COLUMN articles.text_content IS 'The text content of the article, used for search and indexing.';
COMMENT ON COLUMN articles.content IS 'The rich text content of the article, used for display in the UI.';