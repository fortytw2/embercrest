
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
INSERT INTO tiles (name, resistance, dodgebonus, defensebonus) VALUES
                  ('dirt', 1, 0, 0);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DELETE FROM tiles WHERE name = 'dirt';
