-- +goose Up
CREATE TABLE IF NOT EXISTS urls
(
    id    SERIAL  PRIMARY KEY,
    alias TEXT NOT NULL,
    url   TEXT NOT NULL,
    res_url   TEXT NOT NULL,
    UNIQUE (url)
);
CREATE INDEX IF NOT EXISTS idx_alias on urls(alias);
CREATE INDEX IF NOT EXISTS idx_url on urls(alias);

-- +goose Down
drop table urls;