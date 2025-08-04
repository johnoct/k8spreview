package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v3"

	"k8spreview/pkg/k8s"
)

type view int

const (
	listView view = iota
	detailView
	graphView
)

// Model represents the UI state
type Model struct {
	resources []k8s.Resource
	list      list.Model
	selected  *k8s.Resource
	view      view
	viewport  viewport.Model
	width     int
	height    int
}

// Item represents a list item in the UI
type item struct {
	title       string
	description string
	resource    k8s.Resource
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.resource.Kind + "/" + i.resource.Metadata.Name }

// NewModel creates a new UI model
func NewModel(resources []k8s.Resource) Model {
	items := make([]list.Item, len(resources))
	for i, res := range resources {
		items[i] = item{
			title:       fmt.Sprintf("%s/%s", res.Kind, res.Metadata.Name),
                        description: fmt.Sprintf("API Version: %s, Namespace: %s", res.APIVersion, res.Metadata.Namespace),
			resource:    res,
		}
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.Copy().Foreground(lipgloss.NoColor{})
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Copy().Foreground(lipgloss.Color("white"))

	l := list.New(items, delegate, 0, 0)
	l.Title = "Kubernetes Resources"
	l.Styles.Title = TitleStyle
	l.FilterInput.Prompt = "Filter: "
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.FilterPrompt = lipgloss.NewStyle().Foreground(lipgloss.Color("white"))
	l.Styles.FilterCursor = lipgloss.NewStyle().Foreground(lipgloss.Color("white"))

	vp := viewport.New(0, 0)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)

	return Model{
		resources: resources,
		list:      l,
		viewport:  vp,
		view:      listView,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles UI events
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			if m.view == detailView || m.view == graphView {
				m.view = listView
				return m, nil
			}
			return m, tea.Quit
		case "enter":
			if m.view == listView {
				if i, ok := m.list.SelectedItem().(item); ok {
					m.selected = &i.resource
					m.view = detailView

					// Create detailed view with relationships
					yamlData, _ := yaml.Marshal(m.selected)
					content := string(yamlData)

					// Add relationships section
					relations := m.selected.FindRelatedResources(m.resources)
					if len(relations) > 0 {
						content += "\n\nRelationships:\n"
						for _, rel := range relations {
							content += RelationshipStyle.Render(fmt.Sprintf("  %s\n", rel))
						}
					}

					m.viewport.SetContent(content)
				}
			}
		case "g":
			if m.view == listView {
				m.view = graphView
			}
		case "/":
			m.list.ShowFilter()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width, msg.Height)
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height
	}

	if m.view == listView {
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	} else if m.view == detailView {
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}

	return m, cmd
}

// View renders the UI
func (m Model) View() string {
	switch m.view {
	case listView:
		return "\n" + m.list.View()
	case detailView:
		return m.viewport.View()
	case graphView:
		return m.generateGraph()
	default:
		return ""
	}
}

type relationship struct {
	from     string
	to       string
	relation string
}

// generateGraph creates a visual representation of resource relationships
func (m Model) generateGraph() string {
	var relationships []relationship
	resourcesByType := make(map[string][]string)

	// Collect all relationships and organize resources by type
	for _, resource := range m.resources {
		nodeKey := fmt.Sprintf("%s/%s", resource.Kind, resource.Metadata.Name)
		resourcesByType[resource.Kind] = append(resourcesByType[resource.Kind], nodeKey)

		relations := resource.FindRelatedResources(m.resources)
		for _, rel := range relations {
			var r relationship
			if strings.HasPrefix(rel, "→") {
				parts := strings.Split(strings.TrimPrefix(rel, "→ "), "/")
				if len(parts) == 2 {
					r = relationship{
						from:     nodeKey,
						to:       fmt.Sprintf("%s/%s", parts[0], strings.TrimSpace(parts[1])),
						relation: "uses",
					}
				}
			} else if strings.HasPrefix(rel, "←") {
				parts := strings.Split(strings.TrimPrefix(rel, "← Selected by "), "/")
				if len(parts) == 2 {
					r = relationship{
						from:     fmt.Sprintf("%s/%s", parts[0], strings.TrimSpace(parts[1])),
						to:       nodeKey,
						relation: "selects",
					}
				}
			}
			if r.from != "" && r.to != "" {
				relationships = append(relationships, r)
			}
		}
	}

	var sb strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#333333")).
		Padding(0, 1).
		MarginBottom(1)

	sb.WriteString("\n" + titleStyle.Render("Kubernetes Resource Graph") + "\n\n")

	// Legend
	sb.WriteString("Legend:\n")
	for kind, style := range ResourceStyles {
		if len(resourcesByType[kind]) > 0 {
			sb.WriteString(fmt.Sprintf("  %s: %s\n",
				style.Render(fmt.Sprintf("%-10s", kind)),
				style.Render("●")))
		}
	}
	sb.WriteString("  Relationships:\n")
	sb.WriteString(fmt.Sprintf("    %s: Resource uses/references\n",
		GraphEdgeStyle.Render("A ────uses────> B")))
	sb.WriteString(fmt.Sprintf("    %s: Resource selects/targets\n",
		GraphEdgeStyle.Render("A <──selects──── B")))
	sb.WriteString("\n")

	// Resources by type
	sb.WriteString("Resources:\n")
	for kind := range ResourceStyles {
		resources := resourcesByType[kind]
		if len(resources) > 0 {
			style := ResourceStyles[kind]
			sb.WriteString(fmt.Sprintf("  %s\n", style.Render(kind)))
			for _, res := range resources {
				name := strings.Split(res, "/")[1]
				sb.WriteString(fmt.Sprintf("    %s %s\n",
					style.Render("●"),
					style.Render(name)))
			}
		}
	}

	// Relationships
	if len(relationships) > 0 {
		sb.WriteString("\nConnections:\n")
		for _, rel := range relationships {
			fromParts := strings.Split(rel.from, "/")
			toParts := strings.Split(rel.to, "/")

			fromStyle := ResourceStyles[fromParts[0]]
			toStyle := ResourceStyles[toParts[0]]

			var arrow string
			if rel.relation == "uses" {
				arrow = "────uses────>"
			} else {
				arrow = "<──selects────"
			}

			line := fmt.Sprintf("  %s %s %s",
				fromStyle.Render(fromParts[1]),
				GraphEdgeStyle.Render(arrow),
				toStyle.Render(toParts[1]))
			sb.WriteString(line + "\n")
		}
	}

	sb.WriteString("\nPress 'q' to return to list view")
	return sb.String()
}
