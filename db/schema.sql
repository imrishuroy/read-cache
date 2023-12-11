CREATE TABLE "caches" (
  "id" bigserial PRIMARY KEY,
  "title" varchar NOT NULL,
  "link" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);
