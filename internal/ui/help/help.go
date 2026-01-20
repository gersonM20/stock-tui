package help

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Binding struct {
	Key  string
	Desc string
}

type Model struct {
	bindings []Binding
	visible  bool
	width    int
	height   int
}

func New() Model {
	return Model{
		bindings: []Binding{
			{"j/↓", "Move down"},
			{"k/↑", "Move up"},
			{"/", "Search symbols"},
			{"s", "Cycle sort (Name/Price/%)"},
			{"S", "Toggle sort direction"},
			{"Tab", "Cycle time range"},
			{"1-4", "Select time range"},
			{"c", "Cycle chart type"},
			{"r", "Refresh data"},
			{"?", "Toggle help"},
			{"q", "Quit"},
		},
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "?" || msg.String() == "esc" || msg.String() == "q" {
			m.visible = false
		}
	}
	return m, nil
}

func (m *Model) SetSize(w, h int) {
	m.width = w
	m.height = h
}

func (m *Model) Toggle() {
	m.visible = !m.visible
}

func (m *Model) Show() {
	m.visible = true
}

func (m *Model) Hide() {
	m.visible = false
}

func (m Model) Visible() bool {
	return m.visible
}

func (m Model) View() string {
	if !m.visible {
		return ""
	}

	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true).
		Width(10)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#CCCCCC"))

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true).
		MarginBottom(1)

	var sb strings.Builder
	sb.WriteString(titleStyle.Render("Keyboard Shortcuts"))
	sb.WriteString("\n\n")

	for _, b := range m.bindings {
		sb.WriteString(keyStyle.Render(b.Key))
		sb.WriteString(descStyle.Render(b.Desc))
		sb.WriteString("\n")
	}

	content := sb.String()

	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		Background(lipgloss.Color("#1a1a2e"))

	modal := modalStyle.Render(content)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, modal)
}
