package tui

import (
	"errors"
	internalModel "golang_gh/internal/model"

	tea "charm.land/bubbletea/v2"
	"charm.land/glamour/v2"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/tree"
	"github.com/cli/go-gh/v2/pkg/api"
)

func (m *model) renderChildrenView(githubPR *internalModel.GithubPR, currentIndex int) (*tree.Tree, int) {
	var root *tree.Tree
	if m.tree.selectedPRID == githubPR.ID {
		itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
		root = tree.New().Root(itemStyle.Render(githubPR.Title))
	} else {
		root = tree.New().Root(githubPR.Title)
	}

	currentIndex = currentIndex + 1

	returnIndex := 0
	lastIndex := 0
	for i, childPR := range githubPR.Children {
		if len(childPR.Children) > 0 {
			child, childIndex := m.renderChildrenView(childPR, i+currentIndex+returnIndex)
			returnIndex += childIndex
			root.Child(child)
		} else {
			if m.tree.selectedPRID == childPR.ID {
				itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
				root.Child(itemStyle.Render(childPR.Title))
			} else {
				root.Child(childPR.Title)
			}
		}
		lastIndex = i
	}

	return root, returnIndex + lastIndex + 1

}

func humanize(err error) string {
	if httpError, ok := errors.AsType[*api.HTTPError](err); ok {
		switch httpError.StatusCode {
		case 401:
			return "認証が必要です。gh auth login を実行してください"
		case 404:
			return "リポジトリが見つかりません。owner / name を確認してください"
		case 403:
			return "アクセス権がないか、レート制限に達しています"
		}
	}
	return err.Error()
}

func (m model) ViewSelect() tea.View {
	if m.isLoading {
		view := tea.NewView("loading.......")
		view.AltScreen = true
		return view
	}

	// if m.form.State == huh.StateCompleted {
	// 	repositoryOwner := m.form.GetString("RepositoryOwner")
	// 	repositoryName := m.form.GetString("RepositoryName")
	// 	newForm := newRepoForm(repositoryOwner, repositoryName)
	// 	m.form = newForm
	// }

	// formView := m.form.View()

	if m.form.viewMode == SettingViewModeHistory {
		help := m.form.selectForm.Help().ShortHelpView(append(m.form.selectForm.KeyBinds(), registerNewRepositoryKey))
		selectView := m.form.selectForm.View()

		// if m.err != nil {
		// 	formView += "\n" + humanize(m.err)
		// }

		// stack := lipgloss.JoinVertical(lipgloss.Left, selectView+"\n\n", formView)

		view := tea.NewView(selectView + "\n\n" + help)

		view.AltScreen = true
		return view
	}

	registerForm := m.form.registerForm.View()
	view := tea.NewView(registerForm)
	view.AltScreen = true
	return view

}

func (m model) View() tea.View {
	if m.viewMode == ViewModeSetting {
		return m.ViewSelect()
	}

	myTree := tree.Root(".")

	childIndex := 0

	switch m.tree.viewMode {
	case TreeViewModeList:
		{

			for i, githubPR := range m.tree.rootGithubPRList {
				if len(githubPR.Children) == 0 {
					if githubPR.ID == m.tree.selectedPRID {
						itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
						myTree.Child(itemStyle.Render(githubPR.Title))
					} else {
						myTree.Child(githubPR.Title)
					}
				} else {
					child, childIndex2 := m.renderChildrenView(githubPR, i+childIndex)
					childIndex += childIndex2
					myTree.Child(child)
				}
			}

			myTree.Enumerator(tree.RoundedEnumerator)

			// body := lipgloss.JoinVertical(lipgloss.Top, header, myTree.String())

			view := tea.NewView(myTree.String())
			view.AltScreen = true

			// Send the UI for rendering
			return view
		}
	case TreeViewModeDetail:
		{
			githubPR := m.tree.githubPRs[m.tree.selectedPRID]
			r, err := glamour.NewTermRenderer(
				glamour.WithStandardStyle("dark"), // "dark" / "light" / "notty" / "dra/ "ascii" 等
			)

			if err != nil {
				return tea.NewView("")
			}

			out, err := r.Render(githubPR.Body)

			if err != nil {
				return tea.NewView("")
			}

			view := tea.NewView(out)
			view.AltScreen = true

			return view
		}
	default:
		return tea.NewView("")
	}

}
