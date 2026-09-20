package tui

import (
	"fmt"
	internalModel "golang_gh/internal/model"

	tea "charm.land/bubbletea/v2"
	"charm.land/glamour/v2"
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/tree"
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

var organizationNameValue string

func (m model) View() tea.View {

	if m.viewMode == ViewModeSelect {
		if m.form.State == huh.StateCompleted {
			view := tea.NewView(fmt.Sprintf("Github Organization Name: %s", organizationNameValue))
			view.AltScreen = true
			return view
		}

		view := tea.NewView(m.form.View())
		view.AltScreen = true
		return view
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
