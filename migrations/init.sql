ALTER SYSTEM SET max_connections = 1000;
ALTER SYSTEM SET shared_buffers TO "375MB";
ALTER DATABASE postgres SET synchronous_commit=OFF;

CREATE EXTENSION pg_trgm;

CREATE TABLE IF NOT EXISTS pessoas
(
    id         VARCHAR(36) PRIMARY KEY,
    apelido    VARCHAR(32) UNIQUE,
    nome       VARCHAR(100) NOT NULL,
    nascimento CHAR(10)     NOT NULL,
    stack      VARCHAR(1024),
    trigram    TEXT GENERATED ALWAYS AS (
                   LOWER(apelido) || LOWER(nome) || LOWER(stack)
                   ) STORED
);

CREATE INDEX vector_idx ON pessoas USING gist (trigram GIST_TRGM_OPS);
