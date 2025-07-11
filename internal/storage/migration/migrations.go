package migration

import "github.com/Okenamay/shorturl.git/internal/config"

type MigrationEntry struct {
	ID      string
	UpSQL   string
	DownSQL string
}

var migrations = map[string]MigrationEntry{
	"20250520160000": {
		ID: "20250520160000",
		UpSQL: `
            CREATE TABLE IF NOT EXISTS public.urls (
                id BIGSERIAL PRIMARY KEY,
                url VARCHAR(1024),
                short_id VARCHAR(10)
            );
            CREATE UNIQUE INDEX IF NOT EXISTS idx_urls ON public.urls (url, short_id);
        `,
		DownSQL: `
            DROP INDEX IF EXISTS public.idx_urls;
            DROP TABLE IF EXISTS public.urls;
        `,
	},

	// --- БУДУЩИЕ МИГРАЦИИ ДОБАВЛЯТЬ СЮДА ---
	// Пример добавления миграции:
	/*
	   20250630170000: {
	       ID: 20250630170000,
	       Up: `
	           ALTER TABLE public.urls ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
	       `,
	       Down: `
	           ALTER TABLE public.urls DROP COLUMN IF EXISTS created_at;
	       `,
	   },
	*/
}

func DeliverMigration(conf *config.Cfg) MigrationEntry {
	migration := migrations[conf.MigrateID]

	return migration
}
