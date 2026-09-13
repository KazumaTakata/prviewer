package model

import "time"

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
