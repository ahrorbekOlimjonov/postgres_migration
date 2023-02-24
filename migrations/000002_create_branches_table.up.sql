CREATE TABLE IF NOT EXISTS "branches"(
    "id" SERIAL PRIMARY KEY,
    "name" TEXT NOT NULL,
    "phonenumber" VARCHAR[] NOT NULL,
    "store_id" INTEGER REFERENCES  stores(id)
);