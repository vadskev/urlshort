-- +goose Up
CREATE TABLE IF NOT EXISTS urls
(
    id    SERIAL  PRIMARY KEY,
    alias TEXT NOT NULL,
    url   TEXT NOT NULL,
    res_url   TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_alias on urls(alias);

-- +goose Down
drop table note;