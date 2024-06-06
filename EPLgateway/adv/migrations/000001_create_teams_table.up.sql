CREATE TABLE IF NOT EXISTS teams (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    location text NOT NULL,
    stadium text NOT NULL,
    history text NOT NULL
);