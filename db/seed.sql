-- create a table to track posts
CREATE TABLE posts (
    id UUID PRIMARY KEY,
    status TEXT NOT NULL,
    title TEXT NOT NULL,
    summary TEXT NOT NULL,
    content TEXT NOT NULL,
    author_id UUID NOT NULL,
    created_ts timestamp NOT NULL,
    published_ts timestamp,
    updated_ts timestamp NOT NULL
);