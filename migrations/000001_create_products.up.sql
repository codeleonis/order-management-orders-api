CREATE TABLE IF NOT EXISTS products (
    id          UUID         PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    description VARCHAR(500) NOT NULL DEFAULT '',
    price       DECIMAL(10, 2) NOT NULL,
    sku         VARCHAR(50)  NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
