CREATE TABLE cover(
    id SERIAL PRIMARY KEY
);

CREATE TABLE asset(
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    path TEXT NOT NULL UNIQUE,
    kind TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE embedding_jina_clip_v2(
    asset_id UUID PRIMARY KEY REFERENCES asset(id) ON DELETE CASCADE,
    embedding vector(256) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
