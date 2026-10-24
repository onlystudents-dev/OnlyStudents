-- +goose Up
CREATE TABLE files (
    id BIGSERIAL PRIMARY KEY,
    account_uuid UUID,
    type VARCHAR(128) NOT NULL,
    size BIGINT NOT NULL,
    original_name VARCHAR(512),
    uuid UUID NOT NULL,
    created_at DATE NOT NULL
);

-- +goose Down
DROP TABLE files;
