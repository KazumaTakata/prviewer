package tui

import (
	"fmt"
	"log/slog"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"

	internalModel "golang_gh/internal/model"

	"charm.land/huh/v2"
)

type ViewMode int

const (
	ViewModeTree ViewMode = iota
	ViewModeSetting
)

type SettingViewMode int

const (
	SettingViewModeHistory SettingViewMode = iota
	SettingViewModeNew
)

type SettingForm struct {
	selectForm   *huh.Form
	registerForm *huh.Form
	viewMode     SettingViewMode
}

type model struct {
	tree       treeModel
	form       SettingForm
	viewMode   ViewMode
	isLoading  bool
	err        error
	repository internalModel.SettingRepository
}

type treeModel struct {
	cursor           int
	selectedPRID     uint64
	githubPRs        map[uint64]*internalModel.GithubPR
	rootGithubPRList []*internalModel.GithubPR
	viewMode         TreeViewMode
	repositoryOwner  string
	repositoryName   string
	viewport         viewport.Model
}

type TreeViewMode int

const (
	TreeViewModeList TreeViewMode = iota
	TreeViewModeDetail
)

func getSelectForm(repo internalModel.SettingRepository) *huh.Form {

	settingHistories, err := repo.GetRepositorySetting()

	selectForm := huh.NewSelect[internalModel.RepositorySetting]().Title("リポジトリー履歴").Key(SelectFormKey)

	if err != nil {
		slog.Warn("リポジトリ履歴の読み込みに失敗しました", "err", err)
	}

	options := make([]huh.Option[internalModel.RepositorySetting], len(settingHistories))

	for i, settingHistory := range settingHistories {
		option := fmt.Sprintf("%s/%s", settingHistory.RepositoryOwner, settingHistory.RepositoryName)
		options[i] = huh.NewOption(option, settingHistory)
	}

	selectForm.Options(options...)

	selectFormGroup := huh.NewForm(
		huh.NewGroup(
			selectForm,
		),
	).WithShowHelp(false)

	return selectFormGroup
}

func InitializeModel(repo internalModel.SettingRepository) model {

	selectFormGroup := getSelectForm(repo)

	newInputFormGroup := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().Title("新しく入力する"),
			huh.NewInput().
				Title("Github Repository Owner").
				Key("RepositoryOwner").
				Prompt("> "),
			huh.NewInput().
				Title("Github Repository Name").
				Key("RepositoryName").
				Prompt("> "),
		),
	)

	return model{
		tree: treeModel{
			cursor:           0,
			githubPRs:        make(map[uint64]*internalModel.GithubPR, 0),
			rootGithubPRList: make([]*internalModel.GithubPR, 0),
			viewMode:         TreeViewModeList,
			repositoryOwner:  "",
			repositoryName:   "",
			selectedPRID:     0,
			viewport:         viewport.New(),
		},

		form:       SettingForm{selectForm: selectFormGroup, registerForm: newInputFormGroup, viewMode: SettingViewModeHistory},
		viewMode:   ViewModeSetting,
		repository: repo,
		isLoading:  false,
		err:        nil,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return m.form.selectForm.Init()
}
