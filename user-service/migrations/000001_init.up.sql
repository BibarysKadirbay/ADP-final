CREATE TABLE IF NOT EXISTS users (
     id UUID PRIMARY KEY,
     name TEXT NOT NULL,
     email TEXT NOT NULL,
     phone TEXT,
     address TEXT,
     created_at TIMESTAMP DEFAULT NOW()
);