package storage

import (
	"context"
	"database/sql"
	"fmt"
	"hash/fnv"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const migrationLockKey = "aiagent_migrations_v1"

// ApplyMigrations runs SQL files from the migrations directory in a safe, repeatable way.
// It uses a Postgres advisory lock to avoid concurrent runners.
func ApplyMigrations(db *sql.DB) error {
	dir := resolveMigrationsDir()
	if dir == "" {
		return fmt.Errorf("migrations directory not found")
	}

	files, err := filepath.Glob(filepath.Join(dir, "*.sql"))
	if err != nil {
		return fmt.Errorf("list migrations: %w", err)
	}
	sort.Strings(files)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	lockKey := advisoryKey(migrationLockKey)
	if _, err := db.ExecContext(ctx, `SELECT pg_advisory_lock($1)`, lockKey); err != nil {
		return fmt.Errorf("acquire migration lock: %w", err)
	}
	defer func() { _, _ = db.Exec(`SELECT pg_advisory_unlock($1)`, lockKey) }()

	if _, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			filename TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);`); err != nil {
		return fmt.Errorf("ensure schema_migrations: %w", err)
	}

	for _, path := range files {
		name := filepath.Base(path)
		var exists bool
		if err := db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE filename=$1)`, name).Scan(&exists); err != nil {
			return fmt.Errorf("check migration %s: %w", name, err)
		}
		if exists {
			continue
		}
		body, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}
		if strings.TrimSpace(string(body)) == "" {
			continue
		}
		if _, err := db.ExecContext(ctx, string(body)); err != nil {
			return fmt.Errorf("apply migration %s: %w", name, err)
		}
		if _, err := db.ExecContext(ctx, `INSERT INTO schema_migrations (filename) VALUES ($1)`, name); err != nil {
			return fmt.Errorf("record migration %s: %w", name, err)
		}
	}
	return nil
}

func resolveMigrationsDir() string {
	if env := strings.TrimSpace(os.Getenv("MIGRATIONS_DIR")); env != "" {
		if isDir(env) {
			return env
		}
	}
	candidates := []string{
		filepath.Join("migrations"),
		filepath.Join("api", "migrations"),
		filepath.Join("..", "api", "migrations"),
	}
	for _, candidate := range candidates {
		if isDir(candidate) {
			return candidate
		}
	}
	return ""
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func advisoryKey(s string) int64 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return int64(h.Sum32())
}
