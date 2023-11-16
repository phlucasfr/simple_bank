CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "full_name" varchar NOT NULL,
  "cpf_cnpj" varchar UNIQUE NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "is_merchant" boolean DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "last_updated" timestamptz
);

CREATE TABLE "wallets" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "balance" bigint NOT NULL DEFAULT 0,
  "currency" varchar NOT NULL,
  "created_at" timestamptz DEFAULT (now()),
  "country_code" int
);

CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "wallet_id" bigint NOT NULL,
  "amount" bigint NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "transfers" (
  "id" bigserial PRIMARY KEY,
  "from_wallet_id" bigint NOT NULL,
  "to_wallet_id" bigint NOT NULL,
  "amount" bigint NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "wallets" ("owner");

CREATE UNIQUE INDEX ON "wallets" ("owner", "currency");

CREATE INDEX ON "entries" ("wallet_id");

CREATE INDEX ON "transfers" ("from_wallet_id");

CREATE INDEX ON "transfers" ("to_wallet_id");

CREATE INDEX ON "transfers" ("from_wallet_id", "to_wallet_id");

COMMENT ON COLUMN "entries"."amount" IS 'can be negative or positive';

COMMENT ON COLUMN "transfers"."amount" IS 'must be positive';

ALTER TABLE "wallets" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

ALTER TABLE "entries" ADD FOREIGN KEY ("wallet_id") REFERENCES "wallets" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("from_wallet_id") REFERENCES "wallets" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("to_wallet_id") REFERENCES "wallets" ("id");
