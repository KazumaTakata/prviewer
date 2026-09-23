package repository

import (
	"context"
	"database/sql"
	"errors"
	"golang_gh/.gen/model"
	. "golang_gh/.gen/table"

	. "github.com/go-jet/jet/v2/sqlite"

	_ "modernc.org/sqlite"

	internalModel "golang_gh/internal/model"
)

func NewSettingRepository(db *sql.DB) (*SettingRepository, error) {
	if db == nil {
		return nil, errors.New("repository: dbは必須です")
	}

	return &SettingRepository{db: db}, nil
}

type SettingRepository struct {
	db *sql.DB
}

func (s *SettingRepository) GetRepositorySetting() ([]internalModel.RepositorySetting, error) {
	stmt := SELECT(
		Settings.RepositoryName, Settings.RepositoryOwner,
	).FROM(
		Settings,
	)

	var settingsFromDB []model.Settings
	if err := stmt.Query(s.db, &settingsFromDB); err != nil {
		return []internalModel.RepositorySetting{}, nil
	}

	settings := make([]internalModel.RepositorySetting, len(settingsFromDB))

	for i, setting := range settingsFromDB {
		settings[i] = internalModel.RepositorySetting{RepositoryOwner: setting.RepositoryOwner, RepositoryName: setting.RepositoryName}
	}
	return settings, nil
}

func (s *SettingRepository) SaveRepositorySetting(repositorySetting internalModel.RepositorySetting) error {
	stmt := Settings.INSERT(Settings.RepositoryName, Settings.RepositoryOwner).VALUES(repositorySetting.RepositoryName, repositorySetting.RepositoryOwner)
	if _, err := stmt.ExecContext(context.Background(), s.db); err != nil {
		return err
	}

	return nil
}
