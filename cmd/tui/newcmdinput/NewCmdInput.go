package newcmdinput

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	errMsg error
)

type Model struct {
	TextInput textinput.Model
	value     string
	err       error
}

func GetModel() Model {
	ti := textinput.New()
	ti.Placeholder = "My New Chatanium Module"
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 30

	return Model{
		TextInput: ti,
		err:       nil,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			if m.TextInput.Value() != "" {
				m.value = m.TextInput.Value()
				return m, tea.Quit
			}
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.TextInput, cmd = m.TextInput.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return fmt.Sprintf(
		"What is your new module name?\n\n%s\n%s",
		m.TextInput.View(),
		"(esc to quit)",
	) + "\n"
}

func (m Model) GetValue() string {
	return m.value
}
