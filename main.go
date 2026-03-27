package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

func theme() *huh.Theme {
	cyan := lipgloss.Color("6")
	t := huh.ThemeBase()
	t.Focused.Base = lipgloss.NewStyle().PaddingLeft(1)
	t.Focused.Title = lipgloss.NewStyle().Bold(true).PaddingLeft(2)
	t.Focused.SelectSelector = lipgloss.NewStyle().SetString("❯ ").Foreground(cyan)
	t.Focused.SelectedOption = lipgloss.NewStyle().Foreground(cyan)
	t.Blurred = t.Focused
	return t
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

	opts := make([]huh.Option[string], len(lines))
	var selected string
	for i, line := range lines {
		branch := strings.SplitN(line, " ", 2)[1]
		opts[i] = huh.NewOption(line, branch)
		if branch == currentBranch {
			selected = branch
		}
	}

	err = huh.NewSelect[string]().
		Title("Branch:").
		Options(opts...).
		Value(&selected).
		WithTheme(theme()).
		Run()
	if err != nil {
		return err
	}

	cmd := exec.Command("git", "switch", selected)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}
