
-- create database repoboost;
-- CREATE USER repoboost WITH PASSWORD 'repoboost';
-- GRANT ALL PRIVILEGES ON DATABASE repoboost TO repoboost;

CREATE EXTENSION IF NOT EXISTS pgcrypto; -- must be superuser

CREATE OR REPLACE FUNCTION triggerSetUpdatedAt()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updatedAt = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS users (
  id serial PRIMARY KEY,
  createdAt TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updatedAt TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  name varchar(30) NOT NULL,
  username varchar(20) UNIQUE NOT NULL,
  password varchar(72) NOT NULL
);

CREATE TRIGGER setUpdatedAt
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE triggerSetUpdatedAt();
INSERT INTO users (name, username, password) VALUES ('Admin', 'admin', crypt('00000', gen_salt('bf')));
