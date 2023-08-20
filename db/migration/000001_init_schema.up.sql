CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "full_name" varchar NOT NULL,
  "cpf_cnpj" varchar UNIQUE NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "is_merchant" boolean DEFAULT false,
  "created_at" timestamptz DEFAULT (now()),
  "last_updated" timestamptz
);

CREATE TABLE "wallets" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint UNIQUE NOT NULL,
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

CREATE INDEX ON "wallets" ("user_id");

CREATE INDEX ON "entries" ("wallet_id");

CREATE INDEX ON "transfers" ("from_wallet_id");

CREATE INDEX ON "transfers" ("to_wallet_id");

CREATE INDEX ON "transfers" ("from_wallet_id", "to_wallet_id");

COMMENT ON COLUMN "entries"."amount" IS 'can be negative or positive';

COMMENT ON COLUMN "transfers"."amount" IS 'must be positive';

ALTER TABLE "wallets" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "entries" ADD FOREIGN KEY ("wallet_id") REFERENCES "wallets" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("from_wallet_id") REFERENCES "wallets" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("to_wallet_id") REFERENCES "wallets" ("id");
