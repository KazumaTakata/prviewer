package tui

import (
	internalModel "golang_gh/internal/model"
	"log"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"github.com/cli/go-gh/v2/pkg/browser"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	if m.viewMode == ViewModeSelect {
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
		}

		switch msg := msg.(type) {

		// Is it a key press?
		case tea.KeyPressMsg:

			// Cool, what was the actual key pressed?
			switch msg.String() {

			case "tab":
				{
					switch m.viewMode {
					case ViewModeSelect:
						m.viewMode = ViewModeTree
					case ViewModeTree:
						m.viewMode = ViewModeSelect
					}
				}

			case "ctrl+c", "q":
				return m, tea.Quit

			}
		}

		return m, cmd
	}

	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyPressMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		case "tab":
			{
				switch m.viewMode {
				case ViewModeSelect:
					m.viewMode = ViewModeTree
				case ViewModeTree:
					m.viewMode = ViewModeSelect
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
