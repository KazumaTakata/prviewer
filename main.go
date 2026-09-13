package main

import (
	"fmt"
	"os"

	internalTUI "golang_gh/internal/tui"

	tea "charm.land/bubbletea/v2"
)

func main() {
	p := tea.NewProgram(internalTUI.InitializeModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
