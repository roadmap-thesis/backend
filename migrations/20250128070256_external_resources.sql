-- +goose Up
CREATE TYPE resource_type AS ENUM(
  'video',
  'article',
  'book'
);

CREATE TABLE IF NOT EXISTS external_resources(
  "id" integer PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
  "topic_id" integer NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
  "title" character varying(255) NOT NULL,
  "url" character varying(255) NOT NULL,
  "type" resource_type NOT NULL,
  "created_at" timestamp without time zone DEFAULT now() NOT NULL,
  "updated_at" timestamp without time zone DEFAULT now() NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS external_resources;

DROP TYPE IF EXISTS resource_type;

