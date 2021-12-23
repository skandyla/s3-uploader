CREATE TABLE IF NOT EXISTS refresh_tokens (
    id SERIAL NOT NULL unique,
    user_id INT REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    token VARCHAR(255) NOT NULL unique, 
    expires_at TIMESTAMP NOT NULL
);
