ALTER TABLE refresh_tokens
ALTER COLUMN created_at TYPE timestamp USING created_at,
ALTER COLUMN expires_at TYPE timestamp USING expires_at,
ALTER COLUMN revoked_at TYPE timestamp USING revoked_at;

ALTER TABLE articles
ALTER COLUMN created_at TYPE timestamp USING created_at,
ALTER COLUMN updated_at TYPE timestamp USING updated_at;

ALTER TABLE users
ALTER COLUMN created_at TYPE timestamp USING created_at,
ALTER COLUMN updated_at TYPE timestamp USING updated_at;
