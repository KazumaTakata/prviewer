package tui

import (
	"fmt"

	internalModel "golang_gh/internal/model"

	"charm.land/lipgloss/v2"
)

var (
	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1).
			Width(48) // 枠線を含めた幅
	selectedCardStyle = cardStyle.BorderForeground(lipgloss.Color("212"))
	titleStyle        = lipgloss.NewStyle().Bold(true)
	metaStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
)

func renderCard(pr *internalModel.GithubPR, selected bool) string {
	style := cardStyle
	if selected {
		style = selectedCardStyle
	}
	body := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render(pr.Title),
		metaStyle.Render(fmt.Sprintf("@%s · %s", "author")),
	)
	return style.Render(body)
}
