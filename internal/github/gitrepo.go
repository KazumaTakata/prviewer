package github

import (
	"fmt"
	"log"
	"os"
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

func GetGithubPRs() (map[uint64]*internalModel.GithubPR, []uint64) {

	client, err := api.DefaultRESTClient()
	if err != nil {
		log.Fatal(err)
	}
	response := []GithubPRJSON{}

	repositoryOwner := os.Getenv("GH_REPO_OWNER")
	repositoryName := os.Getenv("GH_REPO_NAME")

	err = client.Get(fmt.Sprintf("repos/%s/%s/pulls", repositoryOwner, repositoryName), &response)
	if err != nil {
		log.Fatal(err)
	}

	githubPRMap := make(map[uint64]*internalModel.GithubPR, len(response))

	for _, pr := range response {
		githubPRMap[pr.ID] = &internalModel.GithubPR{ID: pr.ID, Body: pr.Body, Title: pr.Title, URL: pr.HTMLURL, Children: []*internalModel.GithubPR{}, Base: pr.Base.Ref, Head: pr.Head.Ref, CreatedAt: pr.CreatedAt, UpdatedAt: pr.UpdatedAt}
	}

	nonRootPRIDs := make([]uint64, 0)

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

	return githubPRMap, nonRootPRIDs
}
