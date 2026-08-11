# Kubernetes Deployment Guide for PDTS-Go

This project follows a microservices architecture. To deploy this project to Kubernetes, you need to create manifest files (YAML) for each component.

## Best Practices for This Project

1. **Namespace**: 
   - Deploy everything into a dedicated namespace (e.g., `pdts-system`).
   
2. **Infrastructure First**:
   - Deploy MongoDB and Kafka before the services.
   - Use `StatefulSets` for MongoDB and Kafka if you need persistent storage.
   - For MongoDB: Create a `Secret` for `MONGO_INITDB_ROOT_PASSWORD`.
   
3. **Services (Microservices)**:
   - Each service folder in `k8s/services/` should contain:
     - `deployment.yaml`: Defines the desired state (replicas, container image, env variables).
     - `service.yaml`: Defines how to access the pods (ClusterIP, NodePort, or LoadBalancer).
     - `configmap.yaml` (optional): For non-sensitive configurations.
     - `secret.yaml` (optional): For sensitive data like `JWT_SECRET`.

4. **Environment Variables**:
   - Map your `.env` variables to Kubernetes `ConfigMaps` and `Secrets`.
   - Use `envFrom` in your Deployment manifests to load all keys from a ConfigMap/Secret.

5. **Networking**:
   - Use DNS names for service-to-service communication. For example, if `auth-service` is the name of your service in K8s, it can be reached via `http://auth-service:8081`.

## How to create your first manifest

Example `deployment.yaml` for `auth-service`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: auth-service
  namespace: pdts-system
spec:
  replicas: 2
  selector:
    matchLabels:
      app: auth-service
  template:
    metadata:
      labels:
        app: auth-service
    spec:
      containers:
      - name: auth-service
        image: <your-docker-image>
        ports:
        - containerPort: 8081
        env:
        - name: MONGODB_URI
          valueFrom:
            secretKeyRef:
              name: mongodb-secrets
              key: uri
```

## Folder Structure Reference
- `k8s/infrastructure/`: Shared dependencies like Databases and Message Brokers.
- `k8s/services/`: Individual microservice manifests.
