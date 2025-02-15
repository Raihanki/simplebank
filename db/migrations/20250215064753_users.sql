-- +goose Up
-- +goose StatementBegin
CREATE TABLE "users" (
    "username" VARCHAR PRIMARY KEY,
    "password" VARCHAR NOT NULL,
    "full_name" VARCHAR NOT NULL,
    "email" VARCHAR UNIQUE NOT NULL,
    "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");
ALTER TABLE "accounts" ADD CONSTRAINT "accounts_owner_currency_idx" UNIQUE ("owner", "currency");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "accounts_owner_currency_idx";
ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "accounts_owner_fkey";
DROP TABLE IF EXISTS "users";
-- +goose StatementEnd
