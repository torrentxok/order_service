
CREATE TABLE orders (
    order_uid TEXT PRIMARY KEY,
    track_number TEXT NOT NULL,
    entry TEXT NOT NULL,
    locale TEXT NOT NULL,
    internal_signature TEXT NOT NULL,
    customer_id TEXT NOT NULL,
    delivery_service TEXT NOT NULL,
    shardkey TEXT NOT NULL,
    sm_id INTEGER NOT NULL,
    date_created TIMESTAMPTZ NOT NULL,
    oof_shard TEXT NOT NULL
);

CREATE TABLE delivery (
    order_uid TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    phone TEXT NOT NULL,
    zip TEXT NOT NULL,
    city TEXT NOT NULL,
    address TEXT NOT NULL,
    region TEXT NOT NULL,
    email TEXT NOT NULL,
    CONSTRAINT fk_delivery_order
        FOREIGN KEY (order_uid)
        REFERENCES orders(order_uid)
        ON DELETE CASCADE
);

CREATE TABLE payment (
    order_uid TEXT PRIMARY KEY,
    transaction TEXT NOT NULL,
    request_id TEXT,
    currency TEXT NOT NULL,
    provider TEXT NOT NULL,
    amount INTEGER NOT NULL,
    payment_dt INTEGER NOT NULL,
    bank TEXT NOT NULL,
    delivery_cost INTEGER NOT NULL,
    goods_total INTEGER NOT NULL,
    custom_fee INTEGER NOT NULL,
    CONSTRAINT fk_payment_order
        FOREIGN KEY (order_uid)
        REFERENCES orders(order_uid)
        ON DELETE CASCADE
);

CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    order_uid TEXT NOT NULL,
    chrt_id INTEGER NOT NULL,
    track_number TEXT NOT NULL,
    price INTEGER NOT NULL,
    rid TEXT NOT NULL,
    name TEXT NOT NULL,
    sale INTEGER NOT NULL,
    size TEXT NOT NULL,
    total_price INTEGER NOT NULL,
    nm_id INTEGER NOT NULL,
    brand TEXT NOT NULL,
    status INTEGER NOT NULL,
    CONSTRAINT fk_items_order
        FOREIGN KEY (order_uid)
        REFERENCES orders(order_uid)
        ON DELETE CASCADE
);
