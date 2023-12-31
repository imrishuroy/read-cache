CREATE TABLE "caches" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "title" varchar NOT NULL,
  "link" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "caches" ("owner");