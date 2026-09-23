package model

import (
	"time"
)

type GithubPR struct {
	ID        uint64
	Title     string
	URL       string
	Base      string
	Head      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Children  []*GithubPR
	Body      string
}

type RepositorySetting struct {
	RepositoryOwner string
	RepositoryName  string
}

type Setting struct {
	repositorySetting []RepositorySetting
}

type SettingRepository interface {
	SaveRepositorySetting(repositorySetting RepositorySetting) error
	GetRepositorySetting() ([]RepositorySetting, error)
}
