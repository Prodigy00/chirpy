-- noinspection SqlNoDataSourceInspectionForFile

-- +goose Up
CREATE TABLE chirps (
                        id UUID,
                        created_at TIMESTAMP NOT NULL,
                        updated_at TIMESTAMP NOT NULL,
                        body TEXT NOT NULL,
                        user_id UUID REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE chirps;