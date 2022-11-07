DROP TABLE IF EXISTS posts, rss_source;

CREATE TABLE posts (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    link TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    pubdate BIGINT
);

CREATE INDEX idx_links
ON posts(link)