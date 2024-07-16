-- +goose Up
-- +goose StatementBegin
create table places
(
    id      uuid NOT NULL primary key,
    name    text,
    story   text,
    date    date,
    images  text[],
    expenses uuid,
    preview text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE places;
-- +goose StatementEnd