-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.urls (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR(36),
    url VARCHAR(1024),
    short_id VARCHAR(10)
	del_flag BOOLEAN DEFAULT false
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_urls ON public.urls (url, short_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS public.idx_urls;
DROP TABLE IF EXISTS public.urls;
-- +goose StatementEnd