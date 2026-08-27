-- +goose Up
CREATE TABLE refresh_tokens(
    token TEXT PRIMARY KEY,
    created_at TIMESTAMP not null,
    updated_at TIMESTAMP not null,
    expires_at TIMESTAMP not null,
    revoked_at TIMESTAMP,
    userID UUID not null,
    CONSTRAINT fk_id FOREIGN KEY (userID) REFERENCES users (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE refresh_tokens;
