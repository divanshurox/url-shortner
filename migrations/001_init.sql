CREATE TABLE IF NOT EXISTS shorturlmappings (
    short_url TEXT PRIMARY KEY,
    original_url TEXT NOT NULL,
    expires_at  TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS index_short_urls ON shorturlmappings (expires_at);