CREATE TABLE passwords (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE, -- when user in users deleted -> then all his pass will be deleted too 
    service_name TEXT NOT NULL,
    secret_enc BYTEA NOT NULL,
    category TEXT,
    created_at TIMESTAMP DEFAULT now(),
    UNIQUE (user_id, service_name) -- every user can have only one pass to a service
);