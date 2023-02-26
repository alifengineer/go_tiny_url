CREATE TABLE IF NOT EXISTS "users" (
    "id" UUID PRIMARY KEY NOT NULL,
    "first_name" VARCHAR(20) NOT NULL,
    "last_name" VARCHAR(20) NOT NULL,
    "phone" VARCHAR,
    "username" VARCHAR ,
    "password" VARCHAR(1000),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "deleted_at" TIMESTAMP,
    UNIQUE ("phone","username","deleted_at")
);

CREATE TABLE IF NOT EXISTS "urls" (
    "id" UUID PRIMARY KEY NOT NULL,
    "short_url" VARCHAR NOT NULL,
    "long_url" VARCHAR NOT NULL,
    "click_count" INTEGER DEFAULT 0,
    "expire_date" TIMESTAMP,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "deleted_at" TIMESTAMP,
    "limit_click" INTEGER DEFAULT NOT NULL,
    "user_id" UUID REFERENCES "users" ("id"),
    UNIQUE ("short_url","deleted_at")
);
