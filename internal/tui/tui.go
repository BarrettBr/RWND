// Package tui provides the RWND terminal UI.
package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	quitting bool
}

func initialModel() model {
	return model{}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	return fmt.Sprintf(
		"RWND TUI scaffold\n\n"+
			"Press q to quit.\n",
	)
}

// Run starts the terminal UI.
func Run() error {
	p := tea.NewProgram(initialModel())
	_, err := p.Run()
	return err
}
