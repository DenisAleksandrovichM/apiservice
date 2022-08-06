-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.users (
    login varchar(15) NOT NULL PRIMARY KEY,
    first_name varchar(20) NOT NULL,
    last_name varchar(25) NOT NULL,
    weight numeric(5, 2) NOT NULL,
    height smallint NOT NULL,
    age smallint NOT NULL
    )
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.users;
-- +goose StatementEnd
