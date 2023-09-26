CREATE TABLE transactions (
    type INT NOT NULL,
    date TIMESTAMPTZ NOT NULL,
    product_description TEXT NOT NULL,
    product_price_cents INT NOT NULL,
    seller_name TEXT NOT NULL,
    seller_type TEXT NOT NULL
);

