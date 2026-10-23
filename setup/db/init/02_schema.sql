CREATE TABLE app_user(
    id          UUID PRIMARY KEY,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE album_cover(
    id                        UUID PRIMARY KEY DEFAULT uuidv7(),
    title                     TEXT,
    year                      INT,
    artist                    TEXT,
    rhythms                   TEXT[] NOT NULL DEFAULT '{}',
    representation_technique  TEXT NOT NULL DEFAULT '',
    predominance              TEXT NOT NULL DEFAULT '',
    colors                    TEXT[] NOT NULL DEFAULT '{}',
    title_syntax              TEXT NOT NULL DEFAULT '',
    published_by              UUID NOT NULL REFERENCES app_user(id),
    created_at                TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE asset(
    id          UUID PRIMARY KEY DEFAULT uuidv7(),
    path        TEXT NOT NULL UNIQUE,
    kind        TEXT NOT NULL,
    uploaded_by UUID NOT NULL REFERENCES app_user(id),
    pending     BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE album_cover_asset(
    album_cover_id  UUID        NOT NULL REFERENCES album_cover(id),
    asset_id        UUID        NOT NULL REFERENCES asset(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY(album_cover_id, asset_id)
);

CREATE TABLE embedding_jina_clip_v2(
    asset_id UUID PRIMARY KEY REFERENCES asset(id) ON DELETE CASCADE,
    embedding vector(256) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
