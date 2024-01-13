CREATE TABLE "tags" (
    "tag_id" SERIAL PRIMARY KEY,
    "tag_name" VARCHAR(255) UNIQUE NOT NULL
);


CREATE TABLE "user_tags" (
    "user_id" varchar,
    "tag_id" INT,
    PRIMARY KEY ("user_id", "tag_id"),
    FOREIGN KEY ("user_id") REFERENCES users("id"),
    FOREIGN KEY ("tag_id") REFERENCES tags("tag_id")
);

CREATE TABLE "cache_tags" (
    "cache_id" bigserial,
    "tag_id" INT,
    PRIMARY KEY (cache_id, tag_id),
    FOREIGN KEY ("cache_id") REFERENCES  caches("id"),
    FOREIGN KEY ("tag_id") REFERENCES tags("tag_id")
);
