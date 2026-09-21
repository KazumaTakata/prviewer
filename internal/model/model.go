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
	repositoryOwner string
	repositoryName  string
}

type Setting struct {
	repositorySetting []RepositorySetting
}

type SettingRepository interface {
	saveRepositorySetting(repositorySetting RepositorySetting)
	getRepositorySetting() RepositorySetting
}
