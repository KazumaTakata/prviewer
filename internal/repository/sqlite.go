package repository

import (
	"database/sql"

	_ "modernc.org/sqlite"

	internalModel "golang_gh/internal/model"
)

type SettingRepository struct {
	db *sql.DB
}

func (s SettingRepository) saveRepositorySetting(repositorySetting internalModel.RepositorySetting) {

}
