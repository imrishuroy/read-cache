CREATE TABLE "users" (
  "id" varchar PRIMARY KEY,
  "email" varchar NOT NULL,
  "name" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "caches" ADD FOREIGN KEY ("owner") REFERENCES "users" ("id");
