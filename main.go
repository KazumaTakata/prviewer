package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"golang_gh/internal/repository"
	internalTUI "golang_gh/internal/tui"

	_ "modernc.org/sqlite"

	tea "charm.land/bubbletea/v2"
)

const schema = `
CREATE TABLE IF NOT EXISTS settings (
      repository_owner        TEXT NOT NULL,
      repository_name  		  TEXT NOT NULL,
      PRIMARY KEY (repository_owner, repository_name)
);
`

// Open は DB を開き、スキーマを用意して返す。
func Open(ctx context.Context, path string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, fmt.Errorf("ディレクトリ作成: %w", err)
	}
	dsn := path + "?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)"
	sqlDB, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("DB を開けません (%s): %w", path, err)
	}
	if _, err := sqlDB.ExecContext(ctx, schema); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("スキーマ作成: %w", err)
	}
	return sqlDB, nil
}

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	db, err := Open(context.Background(), "./db/ghrepo.db")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	repository, err := repository.NewSettingRepository(db)

	if err != nil {
		fmt.Printf("there's been an error: %v", err)
		os.Exit(1)
	}

	p := tea.NewProgram(internalTUI.InitializeModel(repository))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
