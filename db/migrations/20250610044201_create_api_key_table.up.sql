CREATE TABLE IF NOT EXISTS llm_api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    encrypted_key BYTEA NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_llm_api_keys_user FOREIGN KEY (user_id) REFERENCES users(id)
);
COMMENT ON COLUMN llm_api_keys.encrypted_key IS 'An encrypted LLM Provider API key.';