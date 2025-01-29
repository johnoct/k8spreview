package ui

import "github.com/charmbracelet/lipgloss"

var (
	// TitleStyle is used for main titles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF"))

	// RelationshipStyle is used for relationship text
	RelationshipStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#87CEEB")). // Sky blue
				Italic(true)

	// GraphNodeStyle is used for graph nodes
	GraphNodeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#333333")).
			Padding(0, 1)

	// GraphEdgeStyle is used for graph edges
	GraphEdgeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#87CEEB"))

	// ResourceStyles defines color coding for different resource types
	ResourceStyles = map[string]lipgloss.Style{
		"Service":     lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")), // Green
		"Deployment":  lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00FF")), // Magenta
		"Pod":         lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF")), // Cyan
		"ConfigMap":   lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00")), // Yellow
		"Secret":      lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")), // Red
		"StatefulSet": lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500")), // Orange
		"Ingress":     lipgloss.NewStyle().Foreground(lipgloss.Color("#FF69B4")), // Pink
		"Namespace":   lipgloss.NewStyle().Foreground(lipgloss.Color("#9370DB")), // Purple
		"HPA":         lipgloss.NewStyle().Foreground(lipgloss.Color("#98FB98")), // Pale green
	}
)
