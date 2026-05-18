CREATE TABLE personal_access_tokens(
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    token_hash      BYTEA  NOT NULL UNIQUE,
    user_id         UUID         NOT NULL,
    token_name      VARCHAR(100) NOT NULL,
    last_used_at    TIMESTAMPTZ    NULL,
    expires_at      TIMESTAMPTZ    NULL,
    created_at      TIMESTAMPTZ    NOT NULL DEFAULT NOW(),


    
    CONSTRAINT fk_personal_access_tokens_user
     FOREIGN KEY (user_id)  REFERENCES users(id) ON DELETE CASCADE
    
);

CREATE INDEX idx_pat_user_id ON personal_access_tokens(user_id);    