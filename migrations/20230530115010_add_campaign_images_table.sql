-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS campaign_images (
  id BIGSERIAL PRIMARY KEY,
  campaign_id BIGINT REFERENCES "campaigns" (id),
  file_name VARCHAR(225) NOT NULL,
  is_primary INT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CLOCK_TIMESTAMP(),
  updated_at TIMESTAMPTZ
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table if exists campaigns_images;
-- +goose StatementEnd
