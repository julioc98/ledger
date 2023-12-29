CREATE TABLE entries (
    id SERIAL PRIMARY KEY,
    account VARCHAR(255) NOT NULL,
    direction VARCHAR(255) NOT NULL,
    amount BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);