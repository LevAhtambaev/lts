-- +goose Up
-- +goose StatementBegin
create table travel
(
    id          uuid NOT NULL primary key,
    name        text,
    description text,
    date_start  date,
    date_end    date,
    places      uuid[],
    preview     text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE travel;
-- +goose StatementEnd