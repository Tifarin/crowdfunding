-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS transactions (
  id BIGSERIAL PRIMARY KEY,
  campaign_id BIGINT REFERENCES "campaigns" (id),
  user_id BIGINT REFERENCES "users" (id),
  amount INT,
  payment_url VARCHAR(225),
  status VARCHAR(225) NOT NULL,
  code VARCHAR(225) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CLOCK_TIMESTAMP(),
  updated_at TIMESTAMPTZ
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table if exists transactions;
-- +goose StatementEnd
