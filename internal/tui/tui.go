package tui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

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
	tree               treeModel
	settingHistoryForm *huh.Select[string]
	form               *huh.Form
	viewMode           ViewMode
	err                error
	repository         internalModel.SettingRepository
}

type selectModel struct {
}

type treeModel struct {
	cursor           int
	selectedPRID     uint64
	githubPRs        map[uint64]*internalModel.GithubPR
	rootGithubPRList []*internalModel.GithubPR
	viewMode         TreeViewMode
	repositoryOwner  string
	repositoryName   string
}

type TreeViewMode int

const (
	TreeViewModeList TreeViewMode = iota
	TreeViewModeDetail
)

func InitializeModel(repo internalModel.SettingRepository) model {

	settingHistories, err := repo.GetRepositorySetting()

	selectForm := huh.NewSelect[string]()

	if err == nil {
		for _, settingHistories := range settingHistories {
			option := fmt.Sprintf("%s/%s", settingHistories.RepositoryOwner, settingHistories.RepositoryName)
			selectForm.Options(huh.NewOption(option, option))
		}
	}

	return model{
		tree: treeModel{
			cursor:           0,
			githubPRs:        make(map[uint64]*internalModel.GithubPR, 0),
			rootGithubPRList: make([]*internalModel.GithubPR, 0),
			viewMode:         TreeViewModeList,
		},
		settingHistoryForm: huh.NewSelect[string]().
			Title("Pick a country.").
			Options(
				huh.NewOption("United States", "US"),
				huh.NewOption("Germany", "DE"),
				huh.NewOption("Brazil", "BR"),
				huh.NewOption("Canada", "CA"),
			),
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Github Repository Owner").
					Key("RepositoryOwner").
					Prompt("> "),
				huh.NewInput().
					Title("Github Repository Name").
					Key("RepositoryName").
					Prompt("> "),
			),
		),
		viewMode:   ViewModeSelect,
		repository: repo,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return m.form.Init()
}
