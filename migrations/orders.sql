-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders
(
    id varchar(50) primary key,
    order_data jsonb not null,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd