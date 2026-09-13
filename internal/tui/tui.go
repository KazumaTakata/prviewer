package tui

import (
	"log"
	"os"
	"slices"
	"sort"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/tree"
	"github.com/cli/go-gh/v2/pkg/browser"

	internalGithub "golang_gh/internal/github"
	internalModel "golang_gh/internal/model"

	"charm.land/glamour/v2"
)

type model struct {
	choices          []string         // items on the to-do list
	cursor           int              // which to-do list item our cursor is pointing at
	selected         map[int]struct{} // which to-do items are selected
	selectedPRID     uint64
	githubPRs        map[uint64]*internalModel.GithubPR
	rootGithubPRList []*internalModel.GithubPR
	viewMode         ViewMode
}

type ViewMode int

const (
	ViewModeList ViewMode = iota
	ViewModeDetail
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
		cursor:           0,
		selected:         make(map[int]struct{}),
		githubPRs:        githubPRs,
		rootGithubPRList: rootGithubPRList,
		viewMode:         ViewModeList,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m *model) walkChildren(githubPR *internalModel.GithubPR, currentIndex int) (int, bool) {
	if currentIndex == m.cursor {
		m.selectedPRID = githubPR.ID
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
			if i+currentIndex+returnIndex == m.cursor {
				m.selectedPRID = childPR.ID
				return 0, true
			}
		}
		lastIndex = i
	}
	return returnIndex + lastIndex + 1, false
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyPressMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.githubPRs)-1 {
				m.cursor++
			}

		case "right":
			m.viewMode = ViewModeDetail

		case "left":
			m.viewMode = ViewModeList
		}

	}

	childIndex := 0

	for i, githubPR := range m.rootGithubPRList {
		if len(githubPR.Children) == 0 {
			if i+childIndex == m.cursor {
				m.selectedPRID = githubPR.ID
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
			githubPR := m.githubPRs[m.selectedPRID]

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

func (m *model) renderChildrenView(githubPR *internalModel.GithubPR, currentIndex int) (*tree.Tree, int) {
	var root *tree.Tree
	if m.selectedPRID == githubPR.ID {
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
			if m.selectedPRID == childPR.ID {
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

func (m model) View() tea.View {
	myTree := tree.Root(".")

	childIndex := 0

	switch m.viewMode {
	case ViewModeList:
		{
			for i, githubPR := range m.rootGithubPRList {
				if len(githubPR.Children) == 0 {
					if githubPR.ID == m.selectedPRID {
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

			view := tea.NewView(myTree.String())
			view.AltScreen = true

			// Send the UI for rendering
			return view
		}
	case ViewModeDetail:
		{
			githubPR := m.githubPRs[m.selectedPRID]
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
