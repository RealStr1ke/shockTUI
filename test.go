package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type state int

const (
	initializing state = iota
	ready
)

type model struct {
	state         state
	width, height int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.state = ready
	}
	return m, nil
}

func (m model) View() string {
	if m.state == initializing {
		return "Initializing..."
	}
	return fmt.Sprintf("This is a %dx%d size terminal.", m.width, m.height)
}

func main() {
	if err := tea.NewProgram(model{}).Start(); err != nil {
		fmt.Println("uh oh:", err)
		os.Exit(1)
	}
}