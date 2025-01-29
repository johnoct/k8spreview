/*
Package ui provides the terminal user interface for viewing and navigating Kubernetes resources.

The package implements a TUI (Terminal User Interface) using the Bubble Tea framework
and provides several views for interacting with Kubernetes resources:

Views:
  - List View: Shows all resources in a scrollable, filterable list
  - Detail View: Shows YAML representation and relationships of a selected resource
  - Graph View: Visual representation of resource relationships

Features:
  - Color-coded resource types for better visibility
  - Interactive navigation using keyboard
  - Resource filtering
  - Relationship visualization
  - YAML content viewing
  - Resource graph generation

Navigation:
  - Arrow keys: Navigate through resources
  - Enter: View resource details
  - g: View relationship graph
  - /: Filter resources
  - q: Go back/quit

Example Usage:

	// Start the UI with resources from a YAML file
	if err := ui.Run("deployment.yaml"); err != nil {
	    log.Fatal(err)
	}

	// Or start with pre-parsed resources
	resources := []k8s.Resource{...}
	if err := ui.RunWithResources(resources); err != nil {
	    log.Fatal(err)
	}

The UI is built using the following components:
  - Bubble Tea: Main TUI framework
  - Bubbles: Reusable components (list, viewport)
  - Lip Gloss: Terminal styling
*/
package ui
