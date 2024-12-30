-- +goose Up
-- +goose StatementBegin
ALTER TABLE account ADD COLUMN tok_number VARCHAR(15) NOT NULL COMMENT 'Metok号';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE account DROP COLUMN tok_number;
-- +goose StatementEnd
