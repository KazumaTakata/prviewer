package github

import (
	"fmt"
	"log"
	"slices"
	"sort"
	"time"

	internalModel "golang_gh/internal/model"

	"github.com/cli/go-gh/v2/pkg/api"
)

type GithubPRJSON struct {
	ID        uint64                 `json:"id"`
	URL       string                 `json:"url"`
	HTMLURL   string                 `json:"html_url"`
	Title     string                 `json:"title"`
	Base      GithubPRBaseBranchJSON `json:"base"`
	Head      GithubPRBaseBranchJSON `json:"head"`
	UpdatedAt time.Time              `json:"updated_at"`
	CreatedAt time.Time              `json:"created_at"`
	Body      string                 `json:"body"`
}

type GithubPRBaseBranchJSON struct {
	Ref string `json:"ref"`
}

type GithubOptions struct {
	RepositoryOwner string
	RepositoryName  string
}

func GetGithubPRs(githubOptions GithubOptions) (githubPRMap map[uint64]*internalModel.GithubPR, nonRootPRIDs []uint64, err error) {

	client, err := api.DefaultRESTClient()
	if err != nil {
		log.Fatal(err)
	}
	response := []GithubPRJSON{}

	err = client.Get(fmt.Sprintf("repos/%s/%s/pulls", githubOptions.RepositoryOwner, githubOptions.RepositoryName), &response)
	if err != nil {
		return
	}

	githubPRMap = make(map[uint64]*internalModel.GithubPR, len(response))

	for _, pr := range response {
		githubPRMap[pr.ID] = &internalModel.GithubPR{ID: pr.ID, Body: pr.Body, Title: pr.Title, URL: pr.HTMLURL, Children: []*internalModel.GithubPR{}, Base: pr.Base.Ref, Head: pr.Head.Ref, CreatedAt: pr.CreatedAt, UpdatedAt: pr.UpdatedAt}
	}

	nonRootPRIDs = make([]uint64, 0)

	for prID1 := range githubPRMap {
		for prID2 := range githubPRMap {
			if githubPRMap[prID1].Base == githubPRMap[prID2].Head {
				githubPR := githubPRMap[prID2]
				githubPR.Children = append(githubPR.Children, githubPRMap[prID1])
				githubPRMap[prID2] = githubPR
				nonRootPRIDs = append(nonRootPRIDs, prID1)
			}
		}
	}

	return githubPRMap, nonRootPRIDs, nil
}

func LoadGithubPRs(githubOptions GithubOptions) (rootGithubPRList []*internalModel.GithubPR, githubPRs map[uint64]*internalModel.GithubPR, err error) {
	githubPRs, nonRootPRIDs, err := GetGithubPRs(githubOptions)

	if err != nil {
		return
	}

	githubPRList := []*internalModel.GithubPR{}

	for _, githubPR := range githubPRs {
		githubPRList = append(githubPRList, githubPR)
	}

	sort.Slice(githubPRList, func(i, j int) bool {
		return githubPRList[i].CreatedAt.UnixMilli() > githubPRList[j].CreatedAt.UnixMilli()
	})

	rootGithubPRList = make([]*internalModel.GithubPR, 0)

	for _, githubPR := range githubPRList {
		if !slices.Contains(nonRootPRIDs, githubPR.ID) {
			rootGithubPRList = append(rootGithubPRList, githubPR)
		}
	}

	return rootGithubPRList, githubPRs, nil

}
