CREATE TABLE IF NOT EXISTS users (
    id SERIAL NOT NULL unique,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    registered_at TIMESTAMP NOT NULL DEFAULT NOW()
);
