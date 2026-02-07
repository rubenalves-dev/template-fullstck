-- +goose Up
-- +goose StatementBegin
CREATE TABLE menu_definitions (
    id VARCHAR(120) PRIMARY KEY,
    domain VARCHAR(50) NOT NULL,
    label VARCHAR(100) NOT NULL,
    path VARCHAR(200),
    icon VARCHAR(100),
    order_index INTEGER NOT NULL DEFAULT 0,
    parent_id VARCHAR(120),
    permissions TEXT[] NOT NULL DEFAULT '{}',
    visible BOOLEAN NOT NULL DEFAULT true,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE INDEX idx_menu_definitions_domain ON menu_definitions(domain);
CREATE INDEX idx_menu_definitions_parent ON menu_definitions(parent_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE menu_definitions;
-- +goose StatementEnd
