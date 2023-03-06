CREATE TABLE IF NOT EXISTS orders
(
    id varchar(50) primary key,
    order_data jsonb not null,
);