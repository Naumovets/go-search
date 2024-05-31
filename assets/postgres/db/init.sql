CREATE TABLE IF NOT EXISTS website (
    url VARCHAR(511) NOT NULL PRIMARY KEY,
    title VARCHAR(255),
    content TEXT NOT NULL
);

ALTER TABLE website ADD COLUMN content_tsv tsvector;
CREATE INDEX IF NOT EXISTS idx_tsv_content ON website USING gin(content_tsv);