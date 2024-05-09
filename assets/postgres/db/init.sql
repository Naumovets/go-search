CREATE TABLE IF NOT EXISTS website (
    url VARCHAR(511) NOT NULL PRIMARY KEY,
    content TEXT NOT NULL
);

ALTER TABLE website ADD COLUMN text_tsv tsvector;