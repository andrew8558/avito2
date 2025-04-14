-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE pvz(
    id uuid primary key default uuid_generate_v4(),
    city varchar(256) not null,
    registration_date timestamp not null
);
CREATE TABLE receptions(
    id uuid primary key default uuid_generate_v4(),
    date_time timestamp not null,
    pvz_id uuid not null,
    status varchar(256) not null,
    foreign key (pvz_id) references pvz(id)
);
CREATE TABLE products(
    id uuid primary key default uuid_generate_v4(),
    date_time timestamp not null,
    type varchar(256) not null,
    reception_id uuid not null,
    foreign key (reception_id) references receptions(id)
);
CREATE INDEX idx_receptions_pvz_id ON receptions(pvz_id);
CREATE INDEX idx_products_reception_id ON products(reception_id);
CREATE INDEX idx_receptions_date_time ON receptions(date_time);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP EXTENSION IF EXISTS "uuid-ossp";
DROP TABLE pvz;
DROP TABLE receptions;
DROP TABLE products;
-- +goose StatementEnd
