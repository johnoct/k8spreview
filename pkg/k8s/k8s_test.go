package k8s_test

import (
	"os"
	"path/filepath"
	"testing"

	"k8spreview/pkg/k8s"
)

func TestParseFromFile(t *testing.T) {
	// Create a temporary test file
	content := `apiVersion: v1
kind: Service
metadata:
  name: test-service
  namespace: default
spec:
  selector:
    app: test
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment
  namespace: default
spec:
  selector:
    matchLabels:
      app: test
  template:
    metadata:
      labels:
        app: test`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.yaml")
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test parsing
	resources, err := k8s.ParseFromFile(tmpFile)
	if err != nil {
		t.Fatalf("ParseFromFile failed: %v", err)
	}

	// Verify results
	if len(resources) != 2 {
		t.Errorf("Expected 2 resources, got %d", len(resources))
	}

	// Check first resource (Service)
        if resources[0].Kind != "Service" {
                t.Errorf("Expected Service, got %s", resources[0].Kind)
        }
        if resources[0].Metadata.Name != "test-service" {
                t.Errorf("Expected test-service, got %s", resources[0].Metadata.Name)
        }
        if resources[0].Metadata.Namespace != "default" {
                t.Errorf("Expected namespace default, got %s", resources[0].Metadata.Namespace)
        }
        if resources[0].APIVersion != "v1" {
                t.Errorf("Expected API version v1, got %s", resources[0].APIVersion)
        }

	// Check second resource (Deployment)
        if resources[1].Kind != "Deployment" {
                t.Errorf("Expected Deployment, got %s", resources[1].Kind)
        }
        if resources[1].Metadata.Name != "test-deployment" {
                t.Errorf("Expected test-deployment, got %s", resources[1].Metadata.Name)
        }
        if resources[1].Metadata.Namespace != "default" {
                t.Errorf("Expected namespace default, got %s", resources[1].Metadata.Namespace)
        }
        if resources[1].APIVersion != "apps/v1" {
                t.Errorf("Expected API version apps/v1, got %s", resources[1].APIVersion)
        }
}

func TestFindRelatedResources(t *testing.T) {
	// Create test resources
        service := k8s.Resource{
                APIVersion: "v1",
                Kind:       "Service",
		Metadata: k8s.Metadata{
			Name: "test-service",
		},
		Spec: map[string]interface{}{
			"selector": map[string]interface{}{
				"app": "test",
			},
		},
	}

        deployment := k8s.Resource{
                APIVersion: "apps/v1",
                Kind:       "Deployment",
		Metadata: k8s.Metadata{
			Name: "test-deployment",
		},
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

	resources := []k8s.Resource{service, deployment}

	// Test service finding deployment
	relations := service.FindRelatedResources(resources)
	found := false
	for _, rel := range relations {
		if rel == "→ Selects Deployment/test-deployment" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Service should find related Deployment")
	}

	// Test deployment being found by service
	relations = deployment.FindRelatedResources(resources)
	found = false
	for _, rel := range relations {
		if rel == "← Selected by Service/test-service" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Deployment should be found by Service")
	}
}
