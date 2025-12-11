-- create comments table
CREATE TABLE comments (
                          id BIGSERIAL PRIMARY KEY,
                          entity_id BIGINT NOT NULL,
                          user_id BIGINT NOT NULL,
                          text TEXT NOT NULL,
                          created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                          updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
