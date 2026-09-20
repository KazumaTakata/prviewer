package tui

import (
	"slices"
	"sort"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	internalGithub "golang_gh/internal/github"
	internalModel "golang_gh/internal/model"

	"charm.land/huh/v2"
)

var (
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("62"))
	box         = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
	focused     = box.BorderForeground(lipgloss.Color("212"))
	dim         = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

type ViewMode int

const (
	ViewModeTree ViewMode = iota
	ViewModeSelect
)

type model struct {
	tree     treeModel
	form     *huh.Form
	viewMode ViewMode
}

type selectModel struct {
}

type treeModel struct {
	cursor           int
	selectedPRID     uint64
	githubPRs        map[uint64]*internalModel.GithubPR
	rootGithubPRList []*internalModel.GithubPR
	viewMode         TreeViewMode
}

type TreeViewMode int

const (
	TreeViewModeList TreeViewMode = iota
	TreeViewModeDetail
)

func InitializeModel() model {
	githubPRs, nonRootPRIDs := internalGithub.GetGithubPRs()
	githubPRList := []*internalModel.GithubPR{}

	for _, githubPR := range githubPRs {
		githubPRList = append(githubPRList, githubPR)
	}

	sort.Slice(githubPRList, func(i, j int) bool {
		return githubPRList[i].CreatedAt.UnixMilli() > githubPRList[j].CreatedAt.UnixMilli()
	})

	rootGithubPRList := make([]*internalModel.GithubPR, 0)

	for _, githubPR := range githubPRList {
		if !slices.Contains(nonRootPRIDs, githubPR.ID) {
			rootGithubPRList = append(rootGithubPRList, githubPR)
		}
	}

	return model{
		tree: treeModel{
			cursor:           0,
			githubPRs:        githubPRs,
			rootGithubPRList: rootGithubPRList,
			viewMode:         TreeViewModeList,
		},
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Github Organization Name").
					Value(&organizationNameValue).
					Prompt(">"),
			),
		),
		viewMode: ViewModeSelect,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return m.form.Init()
}
