package migration

import (
	_ "github.com/jackc/pgx/v5/stdlib"
)

// func TestMigrateLauncher_Integration(t *testing.T) {
// 	// Этот тест требует наличия Docker
// 	if testing.Short() {
// 		t.Skip("Skipping integration test in short mode")
// 	}

// 	ctx := context.Background()
// 	appLogger := zap.NewNop().Sugar()

// 	// Проверяем доступность Docker перед запуском - это предотвратит панику
// 	// "rootless Docker"
// 	_, err := testcontainers.NewDockerProvider()
// 	if err != nil {
// 		t.Skipf("Skipping integration test: Docker not available or configured correctly: %v", err)
// 	}

// 	// Запускаем контейнер Postgres для теста
// 	pgContainer, err := postgres.RunContainer(ctx,
// 		testcontainers.WithImage("postgres:15-alpine"),
// 		postgres.WithDatabase("testdb"),
// 		postgres.WithUsername("user"),
// 		postgres.WithPassword("pass"),
// 		testcontainers.WithWaitStrategy(
// 			wait.ForLog("database system is ready to accept connections").
// 				WithOccurrence(2).
// 				WithStartupTimeout(5*time.Minute),
// 		),
// 	)
// 	require.NoError(t, err, "Failed to start Postgres container")
// 	defer func() {
// 		if err := pgContainer.Terminate(ctx); err != nil {
// 			t.Fatalf("Failed to terminate Postgres container: %v", err)
// 		}
// 	}()

// 	dsn, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
// 	require.NoError(t, err, "Failed to get connection string from container")

// 	// Тест успешного запуска миграции
// 	t.Run("Successful Migration 'up'", func(t *testing.T) {
// 		conf := &config.Cfg{PostgreDSN: dsn}

// 		// Запускаем наш MigrateLauncher
// 		err := MigrateLauncher(ctx, conf, appLogger)
// 		assert.NoError(t, err, "MigrateLauncher failed to run")

// 		// Проверяем результат в БД
// 		db, err := sql.Open("pgx", dsn)
// 		require.NoError(t, err, "Failed to connect to test DB for verification")
// 		defer db.Close()

// 		// Проверяем, что goose отработал (создал свою таблицу)
// 		var versionID int64
// 		err = db.QueryRowContext(ctx, "SELECT version_id FROM goose_db_version ORDER BY version_id DESC LIMIT 1").Scan(&versionID)
// 		require.NoError(t, err, "goose_db_version table should exist and have entries")
// 		assert.Equal(t, int64(20250720152000), versionID, "Migration version ID should match our file")

// 		// Проверяем, что наша таблица создана корректно
// 		// (Простой SELECT проверит и наличие таблицы, и всех колонок)
// 		_, err = db.ExecContext(ctx, "SELECT id, user_id, url, short_id, del_flag FROM public.urls LIMIT 1")
// 		assert.NoError(t, err, "public.urls table should exist with all columns")
// 	})

// 	// Тест пропуска миграции, если DSN не указан
// 	t.Run("Migration Skipped (No DSN)", func(t *testing.T) {
// 		conf := &config.Cfg{PostgreDSN: ""} // Пустой DSN
// 		err := MigrateLauncher(ctx, conf, appLogger)
// 		assert.NoError(t, err, "MigrateLauncher should not return error when DSN is empty")
// 	})
// }

// Бенчмаркинг интеграционных тестов с Docker не имеет практического смысла в
// данном контексте, так как измеряет производительность Docker и goose, а не
// производительность Go кода
