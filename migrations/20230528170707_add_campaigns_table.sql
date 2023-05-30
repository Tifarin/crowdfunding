-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS campaigns (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT REFERENCES "users" (id),
  name VARCHAR(225) NOT NULL,
  short_description VARCHAR(255),
  description TEXT,
  perks TEXT,
  backer_count INT,
  goal_amount INT,
  current_amount INT,
  slug VARCHAR(255),
  created_at TIMESTAMPTZ NOT NULL DEFAULT CLOCK_TIMESTAMP(),
  updated_at TIMESTAMPTZ
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table if exists campaigns;
-- +goose StatementEnd
