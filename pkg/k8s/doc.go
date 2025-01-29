/*
Package k8s provides functionality for parsing and analyzing Kubernetes YAML resources.

The package includes support for:
  - Parsing multi-document YAML files containing Kubernetes resources
  - Analyzing relationships between resources (e.g., Service-Deployment connections)
  - Finding references between resources (ConfigMaps, Secrets, etc.)
  - Exporting individual resources to YAML files

Resource Types:
  - Services
  - Deployments
  - StatefulSets
  - Pods
  - ConfigMaps
  - Secrets
  - Ingress
  - HorizontalPodAutoscalers

Relationship Detection:
The package can detect various relationships between resources:
  - Service selector matching Deployment/StatefulSet labels
  - ConfigMap and Secret usage in volumes
  - Secret references in environment variables
  - Ingress backend service references
  - HPA scale target references

Example Usage:

	// Parse resources from a YAML file
	resources, err := k8s.ParseFromFile("deployment.yaml")
	if err != nil {
	    log.Fatal(err)
	}

	// Find relationships for each resource
	for _, res := range resources {
	    relations := res.FindRelatedResources(resources)
	    for _, rel := range relations {
	        fmt.Printf("%s/%s: %s\n", res.Kind, res.Metadata.Name, rel)
	    }
	}
*/
package k8s
