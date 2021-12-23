--CREATE DATABASE s3_uploader;
--\c s3_uploader

CREATE TABLE IF NOT EXISTS users (
    id SERIAL NOT NULL unique,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    registered_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id SERIAL NOT NULL unique,
    user_id INT REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    token VARCHAR(255) NOT NULL unique, 
    expires_at TIMESTAMP NOT NULL
);

-- create record
INSERT INTO users(id, name, email, password, registered_at)
    VALUES('1', 'Serge', 'serge@google.com', '********', '2021-11-08 17:05:23.048055');
