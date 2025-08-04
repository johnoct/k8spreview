package ui

import (
	"regexp"
	"strings"
	"testing"

	"k8spreview/pkg/k8s"
)

func TestGenerateGraphServiceSelectsDeployment(t *testing.T) {
	service := k8s.Resource{
		ApiVersion: "v1",
		Kind:       "Service",
		Metadata:   k8s.Metadata{Name: "test-service"},
		Spec: map[string]interface{}{
			"selector": map[string]interface{}{
				"app": "test",
			},
		},
	}

	deployment := k8s.Resource{
		ApiVersion: "apps/v1",
		Kind:       "Deployment",
		Metadata:   k8s.Metadata{Name: "test-deployment"},
		Spec: map[string]interface{}{
			"template": map[string]interface{}{
				"metadata": map[string]interface{}{
					"labels": map[string]interface{}{
						"app": "test",
					},
				},
			},
		},
	}

	m := NewModel([]k8s.Resource{service, deployment})
	graph := m.generateGraph()

	re := regexp.MustCompile("\x1b\\[[0-9;]*m")
	clean := re.ReplaceAllString(graph, "")

	if !strings.Contains(clean, "test-service <──selects──── test-deployment") {
		t.Fatalf("expected selects edge between service and deployment, got:\n%s", clean)
	}
}
