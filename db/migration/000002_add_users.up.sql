CREATE TABLE "users" (
  "id" varchar PRIMARY KEY,
  "email" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
)