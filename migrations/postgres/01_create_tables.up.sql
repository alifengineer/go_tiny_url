CREATE TABLE IF NOT EXISTS "users" (
    "id" UUID PRIMARY KEY NOT NULL,
    "first_name" VARCHAR(20) NOT NULL,
    "last_name" VARCHAR(20) NOT NULL,
    "phone" VARCHAR,
    "username" VARCHAR ,
    "password" VARCHAR(1000),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "deleted_at" INTEGER DEFAULT 0,
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
    "deleted_at" INTEGER DEFAULT 0,
    "user_id" UUID REFERENCES "users" ("id"),
    "limit_click" INTEGER,
    UNIQUE ("short_url","deleted_at")
);
