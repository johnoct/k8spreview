package k8s

import (
	"fmt"
	"io"
	"os"
	"reflect"

	"gopkg.in/yaml.v3"
)

// Metadata represents Kubernetes resource metadata
type Metadata struct {
	Name      string            `yaml:"name"`
	Namespace string            `yaml:"namespace,omitempty"`
	Labels    map[string]string `yaml:"labels,omitempty"`
}

// Resource represents a Kubernetes resource
type Resource struct {
	ApiVersion string      `yaml:"apiVersion"`
	Kind       string      `yaml:"kind"`
	Metadata   Metadata    `yaml:"metadata"`
	Spec       interface{} `yaml:"spec,omitempty"`
	Data       interface{} `yaml:"data,omitempty"`
}

// Parse parses Kubernetes resources from an io.Reader
func Parse(r io.Reader) ([]Resource, error) {
	var resources []Resource
	decoder := yaml.NewDecoder(r)
	for {
		var resource Resource
		if err := decoder.Decode(&resource); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("error decoding YAML: %w", err)
		}
		resources = append(resources, resource)
	}
	return resources, nil
}

// ParseFromFile parses Kubernetes resources from a YAML file
func ParseFromFile(filename string) ([]Resource, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer f.Close()
	return Parse(f)
}

// Export exports a resource to a YAML file
func Export(resource Resource) error {
	filename := fmt.Sprintf("%s-%s.yaml", resource.Kind, resource.Metadata.Name)
	yamlData, err := yaml.Marshal(resource)
	if err != nil {
		return fmt.Errorf("error marshaling resource: %w", err)
	}

	err = os.WriteFile(filename, yamlData, 0644)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}
	return nil
}

// FindRelatedResources finds resources related to this resource
func (r Resource) FindRelatedResources(allResources []Resource) []string {
	var relations []string
	spec := reflect.ValueOf(r.Spec)

	switch r.Kind {
	case "Service":
		if spec.Kind() == reflect.Map {
			selectorMap := spec.Interface().(map[string]interface{})
			if selector, ok := selectorMap["selector"].(map[string]interface{}); ok {
				selectorStringMap := convertToStringMap(selector)
				// Find Deployments, StatefulSets, and Pods matching the selector
				for _, res := range allResources {
					if res.Kind == "Deployment" || res.Kind == "StatefulSet" || res.Kind == "Pod" {
						if labels := getResourceLabels(res); labels != nil {
							if matchLabels(selectorStringMap, labels) {
								relations = append(relations, fmt.Sprintf("→ Selects %s/%s", res.Kind, res.Metadata.Name))
							}
						}
					}
				}
			}
		}

	case "Deployment", "StatefulSet":
		if spec.Kind() == reflect.Map {
			specMap := spec.Interface().(map[string]interface{})
			if template, ok := specMap["template"].(map[string]interface{}); ok {
				if metadata, ok := template["metadata"].(map[string]interface{}); ok {
					if labels, ok := metadata["labels"].(map[string]interface{}); ok {
						// Find Services that select this workload
						for _, res := range allResources {
							if res.Kind == "Service" {
								svcSpec := reflect.ValueOf(res.Spec)
								if svcSpec.Kind() == reflect.Map {
									svcSpecMap := svcSpec.Interface().(map[string]interface{})
									if selector, ok := svcSpecMap["selector"].(map[string]interface{}); ok {
										selectorStringMap := convertToStringMap(selector)
										labelsStringMap := convertToStringMap(labels)
										if matchLabels(selectorStringMap, labelsStringMap) {
											relations = append(relations, fmt.Sprintf("← Selected by Service/%s", res.Metadata.Name))
										}
									}
								}
							}
						}
					}
				}
			}

			// Check for ConfigMap/Secret references in volumes
			if volumes, ok := findVolumes(specMap); ok {
				for _, vol := range volumes {
					if cm, ok := vol["configMap"].(map[string]interface{}); ok {
						if name, ok := cm["name"].(string); ok {
							relations = append(relations, fmt.Sprintf("→ Uses ConfigMap/%s", name))
						}
					}
					if secret, ok := vol["secret"].(map[string]interface{}); ok {
						if name, ok := secret["name"].(string); ok {
							relations = append(relations, fmt.Sprintf("→ Uses Secret/%s", name))
						}
					}
				}
			}

			// Check for Secret references in environment variables
			if template, ok := specMap["template"].(map[string]interface{}); ok {
				if podSpec, ok := template["spec"].(map[string]interface{}); ok {
					if containers, ok := podSpec["containers"].([]interface{}); ok {
						for _, c := range containers {
							if container, ok := c.(map[string]interface{}); ok {
								if envs, ok := container["env"].([]interface{}); ok {
									for _, e := range envs {
										if env, ok := e.(map[string]interface{}); ok {
											if valueFrom, ok := env["valueFrom"].(map[string]interface{}); ok {
												if secretRef, ok := valueFrom["secretKeyRef"].(map[string]interface{}); ok {
													if name, ok := secretRef["name"].(string); ok {
														relations = append(relations, fmt.Sprintf("→ Uses Secret/%s", name))
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}

	case "Ingress":
		if spec.Kind() == reflect.Map {
			specMap := spec.Interface().(map[string]interface{})
			if rules, ok := specMap["rules"].([]interface{}); ok {
				for _, r := range rules {
					if rule, ok := r.(map[string]interface{}); ok {
						if http, ok := rule["http"].(map[string]interface{}); ok {
							if paths, ok := http["paths"].([]interface{}); ok {
								for _, p := range paths {
									if path, ok := p.(map[string]interface{}); ok {
										if backend, ok := path["backend"].(map[string]interface{}); ok {
											if service, ok := backend["service"].(map[string]interface{}); ok {
												if name, ok := service["name"].(string); ok {
													relations = append(relations, fmt.Sprintf("→ Routes to Service/%s", name))
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}

	case "HorizontalPodAutoscaler":
		if spec.Kind() == reflect.Map {
			specMap := spec.Interface().(map[string]interface{})
			if scaleTargetRef, ok := specMap["scaleTargetRef"].(map[string]interface{}); ok {
				if kind, ok := scaleTargetRef["kind"].(string); ok {
					if name, ok := scaleTargetRef["name"].(string); ok {
						relations = append(relations, fmt.Sprintf("→ Scales %s/%s", kind, name))
					}
				}
			}
		}
	}

	return relations
}

// Helper functions

func findVolumes(spec map[string]interface{}) ([]map[string]interface{}, bool) {
	if template, ok := spec["template"].(map[string]interface{}); ok {
		if podSpec, ok := template["spec"].(map[string]interface{}); ok {
			if volumes, ok := podSpec["volumes"].([]interface{}); ok {
				result := make([]map[string]interface{}, len(volumes))
				for i, v := range volumes {
					if volMap, ok := v.(map[string]interface{}); ok {
						result[i] = volMap
					}
				}
				return result, true
			}
		}
	}
	return nil, false
}

func getResourceLabels(res Resource) map[string]string {
	switch res.Kind {
	case "Pod":
		return res.Metadata.Labels
	case "Deployment", "StatefulSet":
		spec := reflect.ValueOf(res.Spec)
		if spec.Kind() == reflect.Map {
			specMap := spec.Interface().(map[string]interface{})
			if template, ok := specMap["template"].(map[string]interface{}); ok {
				if metadata, ok := template["metadata"].(map[string]interface{}); ok {
					if labels, ok := metadata["labels"].(map[string]interface{}); ok {
						return convertToStringMap(labels)
					}
				}
			}
		}
	}
	return nil
}

func matchLabels(selector, labels map[string]string) bool {
	for k, v := range selector {
		if labels[k] != v {
			return false
		}
	}
	return true
}

func convertToStringMap(m map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		if str, ok := v.(string); ok {
			result[k] = str
		}
	}
	return result
}
