package tui

import (
	github "golang_gh/internal/github"
	internalModel "golang_gh/internal/model"
	"log"
	"os"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"github.com/cli/go-gh/v2/pkg/browser"
)

type prsLoadedMsg struct {
	githubPRs     map[uint64]*internalModel.GithubPR
	rootGithubPRs []*internalModel.GithubPR
	options       github.GithubOptions
}

type prsFailedMsg struct{ err error }

func fetchPRs(options github.GithubOptions) tea.Cmd {
	return func() tea.Msg {
		rootGithubPRs, githubPRs, err := github.LoadGithubPRs(options)
		log.Printf("fetchPRs: %v", rootGithubPRs)
		if err != nil {
			return prsFailedMsg{err: err}
		}
		return prsLoadedMsg{githubPRs: githubPRs, rootGithubPRs: rootGithubPRs, options: options}
	}
}

func (m model) UpdateHistorySetting(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyPressMsg, ok := msg.(tea.KeyPressMsg); ok {
		log.Printf("update msg.Key().Text: %s", keyPressMsg.String())
		if keyPressMsg.String() == "enter" {
			value := m.form.selectForm.GetFocusedField().GetValue()
			if repositorySettings, ok := value.(internalModel.RepositorySetting); ok {

				options := github.GithubOptions{
					RepositoryName:  repositorySettings.RepositoryName,
					RepositoryOwner: repositorySettings.RepositoryOwner,
				}
				//
				// 					newForm := newRepoForm(repositorySettings.RepositoryOwner, repositorySettings.RepositoryName)
				// 					m.form = newForm
				m.isLoading = true
				//
				return m, fetchPRs(options)

			}

			log.Printf("update getGetFocusedField: %+v", value)
		}

	}

	form, cmd := m.form.selectForm.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form.selectForm = f
	}

	switch msg := msg.(type) {
	case prsLoadedMsg:
		m.isLoading = false
		m.tree.rootGithubPRList = msg.rootGithubPRs
		m.tree.githubPRs = msg.githubPRs
		log.Printf("update prsLoadedMsg: %v", m.tree.rootGithubPRList)

		m.viewMode = ViewModeTree
		m.repository.SaveRepositorySetting(internalModel.RepositorySetting{
			RepositoryOwner: msg.options.RepositoryOwner,
			RepositoryName:  msg.options.RepositoryName,
		})
		return m, nil

	case prsFailedMsg:
		m.err = msg.err
		return m, nil

	// Is it a key press?
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, registerNewRepositoryKey):
			{
				m.form.viewMode = SettingViewModeNew
				m.form.registerForm =
					huh.NewForm(
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
				m.form.registerForm.Init()

				form, _ = m.form.registerForm.Update(struct{}{})

				if f, ok := form.(*huh.Form); ok {
					m.form.registerForm = f
				}

			}

		}
	}

	return m, cmd

}

func newRepoForm(prevOwner, prevName string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewNote().Title("新しく入力する"),
			huh.NewInput().
				Title("Github Repository Owner").
				Key("RepositoryOwner").
				Prompt("> ").
				Value(&prevOwner), // ← 前回値を入れておく
			huh.NewInput().
				Title("Github Repository Name").
				Key("RepositoryName").
				Prompt("> ").
				Value(&prevName),
		),
	)
}

func (m model) UpdateSetting(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.form.viewMode == SettingViewModeHistory {
		return m.UpdateHistorySetting(msg)
	}

	if m.form.registerForm.State == huh.StateCompleted {
		repositoryOwner := m.form.registerForm.GetString("RepositoryOwner")
		repositoryName := m.form.registerForm.GetString("RepositoryName")

		m.tree.repositoryOwner = repositoryOwner
		m.tree.repositoryName = repositoryName
		options := github.GithubOptions{
			RepositoryName:  repositoryName,
			RepositoryOwner: repositoryOwner,
		}

		newForm := newRepoForm(repositoryOwner, repositoryName)
		m.form.registerForm = newForm
		m.isLoading = true

		return m, fetchPRs(options)

	}

	switch msg := msg.(type) {
	case prsLoadedMsg:
		m.isLoading = false
		m.tree.rootGithubPRList = msg.rootGithubPRs
		m.tree.githubPRs = msg.githubPRs
		log.Printf("update prsLoadedMsg: %v", m.tree.rootGithubPRList)

		m.viewMode = ViewModeTree
		m.repository.SaveRepositorySetting(internalModel.RepositorySetting{
			RepositoryOwner: msg.options.RepositoryOwner,
			RepositoryName:  msg.options.RepositoryName,
		})
		return m, nil

	case prsFailedMsg:
		m.err = msg.err
		m.isLoading = false
		form, cmd := m.form.registerForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form.registerForm = f
		}
		return m, cmd

	// Is it a key press?
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, registerNewRepositoryKey):
			{

			}

		}

		// Cool, what was the actual key pressed?
		switch msg.String() {

		case "tab":
			{
				switch m.form.viewMode {
				case SettingViewModeHistory:
					m.form.viewMode = SettingViewModeNew
					m.form.registerForm.Init()
					form, cmd := m.form.registerForm.Update(msg)
					if f, ok := form.(*huh.Form); ok {
						m.form.registerForm = f
					}

					return m, cmd
				case SettingViewModeNew:
					m.form.viewMode = SettingViewModeHistory
					//
					// selectForm := getSelectForm(m.repository)
					// m.form.selectForm = selectForm
					//
					form, cmd := m.form.selectForm.Update(msg)
					if f, ok := form.(*huh.Form); ok {
						m.form.selectForm = f
					}
					m.form.selectForm.Init()
					return m, cmd
				}
			}

		case "ctrl+c", "q":
			return m, tea.Quit

		}
	}

	form, cmd := m.form.registerForm.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form.registerForm = f
	}

	return m, cmd
}

func (m model) UpdateTree(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyPressMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		case "tab":
			{
				switch m.viewMode {
				case ViewModeSetting:
					m.viewMode = ViewModeTree
				case ViewModeTree:
					m.viewMode = ViewModeSetting
				}
			}

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.tree.cursor > 0 {
				m.tree.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.tree.cursor < len(m.tree.githubPRs)-1 {
				m.tree.cursor++
			}

		case "right":
			m.tree.viewMode = TreeViewModeDetail

		case "left":
			m.tree.viewMode = TreeViewModeList
		}

	}

	childIndex := 0

	for i, githubPR := range m.tree.rootGithubPRList {
		if len(githubPR.Children) == 0 {
			if i+childIndex == m.tree.cursor {
				m.tree.selectedPRID = githubPR.ID
				break
			}
		} else {
			childIndex2, found := m.walkChildren(githubPR, i+childIndex)
			if found {
				break
			}

			childIndex += childIndex2
		}
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter", "space":
			githubPR := m.tree.githubPRs[m.tree.selectedPRID]

			b := browser.New("", os.Stdout, os.Stderr)
			if err := b.Browse(githubPR.URL); err != nil {
				log.Fatal(err)
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil

}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.viewMode == ViewModeSetting {
		return m.UpdateSetting(msg)
	}

	return m.UpdateTree(msg)
}

func (m *model) walkChildren(githubPR *internalModel.GithubPR, currentIndex int) (int, bool) {
	if currentIndex == m.tree.cursor {
		m.tree.selectedPRID = githubPR.ID
		return 0, true
	}

	currentIndex = currentIndex + 1

	returnIndex := 0
	lastIndex := 0
	for i, childPR := range githubPR.Children {
		if len(childPR.Children) > 0 {
			childIndex, found := m.walkChildren(childPR, i+currentIndex+returnIndex)
			if found {
				return 0, true
			}

			returnIndex += childIndex
		} else {
			if i+currentIndex+returnIndex == m.tree.cursor {
				m.tree.selectedPRID = childPR.ID
				return 0, true
			}
		}
		lastIndex = i
	}
	return returnIndex + lastIndex + 1, false
}
