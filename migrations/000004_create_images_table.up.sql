CREATE TABLE images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    image BYTEA NOT NULL
);