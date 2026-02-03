package migration

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Migration struct {
	Version    int64
	VersionRaw string
	Filename   string
	Path       string
	DownPath   string
	Applied    bool
}

var filenameRe = regexp.MustCompile(`^(\d+)_.*\.sql$`)

func LoadDir(ctx context.Context, pool *pgxpool.Pool, dir string) ([]Migration, error) {
	if _, err := pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS schema_versions (version TEXT NOT NULL PRIMARY KEY)`); err != nil {
		return nil, err
	}
	applied := make(map[string]struct{})
	rows, err := pool.Query(ctx, `SELECT version FROM schema_versions`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		applied[v] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	migrations := make([]Migration, 0)
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, ".down.sql") {
			continue
		}
		if !strings.HasSuffix(name, ".sql") {
			continue
		}
		m := filenameRe.FindStringSubmatch(name)
		if m == nil {
			continue
		}
		v, err := strconv.ParseInt(m[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse migration version from %s: %w", name, err)
		}
		_, isApplied := applied[m[1]]
		migrations = append(migrations, Migration{
			Version:    v,
			VersionRaw: m[1],
			Filename:   name,
			Path:       filepath.Join(dir, name),
			DownPath:   filepath.Join(dir, strings.Replace(name, ".sql", ".down.sql", 1)),
			Applied:    isApplied,
		})
	}

	sort.Slice(migrations, func(i, j int) bool {
		if migrations[i].Version == migrations[j].Version {
			return migrations[i].Filename < migrations[j].Filename
		}
		return migrations[i].Version < migrations[j].Version
	})

	for i := 1; i < len(migrations); i++ {
		if migrations[i].Version == migrations[i-1].Version {
			return nil, fmt.Errorf("duplicate migration version %d", migrations[i].Version)
		}
	}

	return migrations, nil
}

func ApplyAll(ctx context.Context, pool *pgxpool.Pool, migrationsDir string) error {
	migrations, err := LoadDir(ctx, pool, migrationsDir)
	if err != nil {
		return err
	}

	for _, m := range migrations {
		if m.Applied {
			continue
		}

		sqlBytes, err := os.ReadFile(m.Path)
		if err != nil {
			return err
		}
		sqlText := strings.TrimSpace(string(sqlBytes))
		if sqlText == "" {
			return fmt.Errorf("migration %s is empty", m.Filename)
		}

		tx, err := pool.Begin(ctx)
		if err != nil {
			return err
		}

		if _, err := tx.Exec(ctx, sqlText); err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("apply %s: %w", m.Filename, err)
		}

		if _, err := tx.Exec(ctx, `INSERT INTO schema_versions (version) VALUES ($1)`, m.VersionRaw); err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("record %s: %w", m.Filename, err)
		}

		if err := tx.Commit(ctx); err != nil {
			return err
		}

	}

	return nil
}

func RollbackLast(ctx context.Context, pool *pgxpool.Pool, migrationsDir string) error {
	migrations, err := LoadDir(ctx, pool, migrationsDir)
	if err != nil {
		return err
	}

	var lastApplied *Migration
	for i := len(migrations) - 1; i >= 0; i-- {
		if migrations[i].Applied {
			lastApplied = &migrations[i]
			break
		}
	}

	if lastApplied == nil {
		fmt.Println("No migrations to rollback.")
		return nil
	}

	sqlBytes, err := os.ReadFile(lastApplied.DownPath)
	if err != nil {
		return fmt.Errorf("reading down migration %s: %w", lastApplied.Filename, err)
	}
	sqlText := strings.TrimSpace(string(sqlBytes))
	if sqlText == "" {
		return fmt.Errorf("down migration %s is empty", lastApplied.Filename)
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, sqlText); err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("rollback %s: %w", lastApplied.Filename, err)
	}

	if _, err := tx.Exec(ctx, `DELETE FROM schema_versions WHERE version = $1`, lastApplied.VersionRaw); err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("unrecord %s: %w", lastApplied.Filename, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	fmt.Printf("Successfully rolled back migration: %s\n", lastApplied.Filename)

	return nil
}
