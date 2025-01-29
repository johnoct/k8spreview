package ui

import (
	"fmt"

	"k8spreview/pkg/k8s"

	tea "github.com/charmbracelet/bubbletea"
)

// Run starts the UI application
func Run(yamlPath string) error {
	resources, err := k8s.ParseFromFile(yamlPath)
	if err != nil {
		return fmt.Errorf("failed to parse YAML file: %w", err)
	}

	p := tea.NewProgram(
		NewModel(resources),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	_, err = p.Run()
	return err
}

// RunWithResources starts the UI application with pre-parsed resources
func RunWithResources(resources []k8s.Resource) error {
	p := tea.NewProgram(
		NewModel(resources),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	_, err := p.Run()
	return err
}
