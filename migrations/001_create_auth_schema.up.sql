CREATE SCHEMA IF NOT EXISTS auth;

CREATE TABLE auth.users (
                            id SERIAL PRIMARY KEY,
                            email VARCHAR(255) NOT NULL UNIQUE,
                            password_hash VARCHAR(255) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_email ON auth.users(email);

CREATE TABLE auth.apps (
                           id INTEGER PRIMARY KEY,
                           name VARCHAR(255) NOT NULL UNIQUE,
                           secret VARCHAR(255) NOT NULL UNIQUE
);
