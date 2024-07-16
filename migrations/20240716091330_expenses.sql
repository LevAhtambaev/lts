-- +goose Up
-- +goose StatementBegin
create table expenses
(
    id            uuid NOT NULL primary key,
    road          integer,
    residence     integer,
    food          integer,
    entertainment integer,
    other         integer
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE expenses;
-- +goose StatementEnd