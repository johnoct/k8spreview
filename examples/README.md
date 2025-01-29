# Example Kubernetes Resources

This directory contains example YAML files demonstrating various Kubernetes resource types and their relationships.

## multi-resource.yaml

This file contains a complete example of a web application deployment with the following resources:

1. **Service** (`web-service`)
   - Exposes port 80
   - Selects pods with label `app: web`

2. **Deployment** (`web-deployment`)
   - Runs 3 replicas of an nginx container
   - Uses ConfigMap for configuration
   - Uses Secret for database password
   - Labels pods with `app: web`

3. **ConfigMap** (`web-config`)
   - Contains JSON configuration
   - Mounted as a volume in the Deployment

4. **Secret** (`db-secret`)
   - Contains base64-encoded password
   - Used as an environment variable in the Deployment

5. **Ingress** (`web-ingress`)
   - Routes HTTP traffic to the Service
   - Uses path-based routing

6. **HorizontalPodAutoscaler** (`web-hpa`)
   - Automatically scales the Deployment
   - Based on CPU utilization
   - Scales between 1 and 10 replicas

### Resource Relationships

The example demonstrates several types of relationships between resources:

- Service → Deployment (via label selector)
- Deployment → ConfigMap (via volume mount)
- Deployment → Secret (via environment variable)
- Ingress → Service (via backend reference)
- HPA → Deployment (via scale target reference)

### Usage

To view these resources and their relationships using k8spreview:

```bash
k8spreview examples/multi-resource.yaml
```

Navigate through the resources using:
- Arrow keys to move through the list
- Enter to view resource details
- 'g' to view the relationship graph
- 'q' to go back or quit
