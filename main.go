package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	cyan     = lipgloss.Color("6")
	selected = lipgloss.NewStyle().Foreground(cyan)
	caret    = lipgloss.NewStyle().SetString("❯ ").Foreground(cyan)
	title    = lipgloss.NewStyle().Bold(true)
	indent   = strings.Repeat(" ", lipgloss.Width(caret.String()))
)

const pageSize = 9

type model struct {
	branches []string
	labels   []string
	cursor   int
	offset   int
	choice   string
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
			m.offset = m.centerOffset()
		case "down", "j":
			if m.cursor < len(m.branches)-1 {
				m.cursor++
			}
			m.offset = m.centerOffset()
		case "enter":
			m.choice = m.branches[m.cursor]
			return m, tea.Quit
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) centerOffset() int {
	offset := m.cursor - pageSize/2
	offset = max(offset, 0)
	offset = min(offset, max(len(m.branches)-pageSize, 0))
	return offset
}

func (m model) View() string {
	var sb strings.Builder
	sb.WriteString(indent + title.Render("Branch:") + "\n")

	end := min(m.offset+pageSize, len(m.labels))
	for i := m.offset; i < end; i++ {
		if i == m.cursor {
			sb.WriteString(caret.String() + selected.Render(m.labels[i]))
		} else {
			sb.WriteString(indent + m.labels[i])
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func run() error {
	out, err := exec.Command(
		"git",
		"for-each-ref",
		"--sort=committerdate",
		"refs/heads/",
		"--format=%(committerdate:short) %(refname:short)",
	).Output()
	if err != nil {
		return err
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	current, _ := exec.Command("git", "branch", "--show-current").Output()
	currentBranch := strings.TrimSpace(string(current))

	branches := make([]string, len(lines))
	labels := make([]string, len(lines))
	cursor := 0
	for i, line := range lines {
		parts := strings.SplitN(line, " ", 2)
		branches[i] = parts[1]
		labels[i] = line
		if parts[1] == currentBranch {
			cursor = i
		}
	}

	offset := max(cursor-pageSize/2, 0)
	if offset+pageSize > len(branches) {
		offset = max(len(branches)-pageSize, 0)
	}

	m, err := tea.NewProgram(model{
		branches: branches,
		labels:   labels,
		cursor:   cursor,
		offset:   offset,
	}).Run()
	if err != nil {
		return err
	}

	choice := m.(model).choice
	if choice == "" {
		return fmt.Errorf("cancelled")
	}

	cmd := exec.Command("git", "switch", choice)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}
