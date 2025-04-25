package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/v2"
)

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	km, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch km.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "up":
		if m.cursor > 0 {
			m.cursor--
		}

	case "down":
		if m.cursor < len(m.choices)-1 {
			m.cursor++
		}

	case "enter", " ":
		_, ok := m.selected[m.cursor]
		if ok {
			delete(m.selected, m.cursor)
		} else {
			m.selected[m.cursor] = struct{}{}

			for k := range m.selected {
				if k != m.cursor {
					delete(m.selected, k)
				}
			}
		}
	}

	return m, nil
}

func (m model) View() string {

	s := "Choose message to send\n"

	for i, choice := range m.choices {

		style := lipgloss.NewStyle()

		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		_, ok := m.selected[i]
		if ok {
			checked = "x"
			style = style.Foreground(lipgloss.Green)
		}

		s += style.Render(fmt.Sprintf("%s [%s] %s", cursor, checked, choice))
		s += "\n"
	}

	s += "\nPress q to quit.\n"

	return s
}
