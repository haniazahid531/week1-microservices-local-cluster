# Week 1 — Microservices & Local Kubernetes Cluster

## Overview
This project contains two small Go microservices, optimized multi-stage Dockerfiles, a Terraform-provisioned local Kubernetes cluster using kind, and raw Kubernetes Deployments and Services.

## Architecture

```text
Client / curl
    |
    v
frontend-service:8081
    |
    v
backend-service:8080
```

The frontend reads `BACKEND_URL=http://backend-service:8080` and calls the backend through Kubernetes service discovery.

## Repository Structure

```text
backend/       Backend source code and Dockerfile
frontend/      Frontend source code and Dockerfile
terraform/     Terraform configuration for the local kind cluster
k8s/           Raw Kubernetes YAML manifests
screenshots/   Evidence captured during testing
```

## Prerequisites
- Linux VM
- Git
- Docker
- kubectl
- kind
- Terraform
- Helm

Record the installed versions here:

```bash
git --version
docker --version
kubectl version --client
kind version
terraform version
helm version 
'''

Installed versions used:

- Git: 2.53.0
- Docker: 28.5.2+dfsg4
- kubectl: v1.36.3
- Kustomize: v5.8.1
- kind: v0.32.0
- Terraform: v1.15.8
- Helm: v4.2.3

## 1. Build the Images

```bash
docker build -t week1-backend:v1 ./backend
docker build -t week1-frontend:v1 ./frontend
docker images | grep week1
```

## 2. Test the Containers Locally

```bash
docker network create week1-net

docker run -d --name backend \
  --network week1-net \
  -p 8080:8080 \
  week1-backend:v1

docker run -d --name frontend \
  --network week1-net \
  -p 8081:8081 \
  -e BACKEND_URL=http://backend:8080 \
  week1-frontend:v1

curl http://localhost:8080/health
curl http://localhost:8081/

docker rm -f frontend backend
docker network rm week1-net
```

## 3. Provision the Cluster with Terraform

```bash
cd terraform
terraform fmt -check
terraform init
terraform validate
terraform plan
terraform apply
export KUBECONFIG="$PWD/kubeconfig"
kubectl cluster-info
kubectl get nodes
cd ..
```

## 4. Load Images into kind

```bash
kind load docker-image week1-backend:v1 --name week1-cluster
kind load docker-image week1-frontend:v1 --name week1-cluster
```

## 5. Deploy Raw Kubernetes Manifests

```bash
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/backend.yaml
kubectl apply -f k8s/frontend.yaml
kubectl wait --for=condition=Available deployment/backend -n week1 --timeout=120s
kubectl wait --for=condition=Available deployment/frontend -n week1 --timeout=120s
kubectl get all -n week1
```

## 6. Verify Internal Communication with kubectl exec and curl

```bash
kubectl run curlbox -n week1 \
  --image=curlimages/curl:8.12.1 \
  --command -- sleep 3600

kubectl wait --for=condition=Ready pod/curlbox -n week1 --timeout=120s
kubectl exec -n week1 curlbox -- curl -s http://backend-service:8080/health
kubectl exec -n week1 curlbox -- curl -s http://frontend-service:8081/
kubectl delete pod curlbox -n week1
```

## 7. Test from the VM

```bash
kubectl port-forward -n week1 service/frontend-service 8081:8081
```

In another terminal:

```bash
curl http://localhost:8081/
```

## Expected Result
The frontend response should contain its own message and the JSON returned by the backend, proving that both services are running and communicating through the Kubernetes Service.

## Troubleshooting

### Docker permission denied
Log out of Kali and log back in after adding your user to the `docker` group. Check with:

```bash
groups
docker ps
```

### ImagePullBackOff
Reload both local images into the named kind cluster and confirm the manifests use versioned tags with `imagePullPolicy: IfNotPresent`.

### kubectl cannot reach the cluster

```bash
cd terraform
export KUBECONFIG="$PWD/kubeconfig"
kubectl get nodes
```

### Pod logs and details

```bash
kubectl get pods -n week1
kubectl describe pod -n week1 <pod-name>
kubectl logs -n week1 deployment/backend
kubectl logs -n week1 deployment/frontend
```

## Cleanup

```bash
kubectl delete namespace week1
cd terraform
terraform destroy
cd ..
```

## Evidence to Include
Add screenshots or copied terminal output for:
1. Prerequisite versions.
2. Successful local Docker test.
3. `terraform apply` and `kubectl get nodes`.
4. `kubectl get all -n week1`.
5. Both `kubectl exec ... curl` commands.
6. Final frontend response.

## What I Built and My Approach
I built two simple Go HTTP microservices: a frontend service and a backend service. I used multi-stage Dockerfiles so the applications could be compiled in one stage and copied into smaller final images. I used kind because it provides a lightweight local Kubernetes cluster through Docker. Terraform was used to create the cluster, while raw Kubernetes Deployment and Service manifests were used to deploy both applications. I verified the setup by checking the pods and services, using kubectl exec with curl for internal communication, and testing the frontend through port forwarding.

---

# Week 2 — Helm, Istio and Strict mTLS

## Overview

In Week 2, the Week 1 raw Kubernetes manifests were converted into a reusable Helm chart. Istio was then installed on the local kind cluster, automatic Envoy sidecar injection was enabled, and strict mutual TLS was enforced between the frontend and backend services.

The final environment includes:

- Frontend and backend managed through Helm
- Istio control plane installed using Helm
- Automatic Envoy sidecar injection
- Strict namespace-wide mTLS
- Successful communication between injected workloads
- Blocked plaintext access from an un-injected pod
- Prometheus metrics collection
- Kiali service-mesh topology visualization

## Architecture

```text
curl-injected pod
        |
        | Istio mTLS
        v
frontend-service:8081
        |
        | Istio mTLS
        v
backend-service:8080

```

## Project Structure

```text
helm/week2-microservices/
├── Chart.yaml
├── values.yaml
└── templates/
    ├── backend-deployment.yaml
    ├── backend-service.yaml
    ├── frontend-deployment.yaml
    └── frontend-service.yaml

istio/
├── strict-mtls.yaml
└── test-pods.yaml

evidence/
├── week2-final-status.txt
└── week2-mtls-test.txt
```

## Helm Deployment

The Week 1 raw Kubernetes manifests were converted into a reusable Helm chart.

Validate and install the chart:

```bash
helm lint helm/week2-microservices

helm install week2-app helm/week2-microservices \
  --namespace week2 \
  --create-namespace \
  --wait \
  --timeout 5m
```

Verify the release:

```bash
helm list -n week2
kubectl get all -n week2
```

## Istio Installation

Istio was installed using the official Helm repository:

```bash
helm repo add istio https://istio-release.storage.googleapis.com/charts
helm repo update
```

The Istio base resources and control plane were installed:

```bash
helm install istio-base istio/base \
  --version 1.30.3 \
  --namespace istio-system \
  --create-namespace \
  --set defaultRevision=default \
  --wait

helm install istiod istio/istiod \
  --version 1.30.3 \
  --namespace istio-system \
  --wait
```

## Automatic Sidecar Injection

Automatic Envoy sidecar injection was enabled:

```bash
kubectl label namespace week2 \
  istio-injection=enabled --overwrite

kubectl rollout restart deployment -n week2
```

The frontend and backend pods then showed `2/2 Running`, representing the application container and Istio Envoy sidecar.

## Strict mTLS

Strict mutual TLS was enabled using:

```yaml
apiVersion: security.istio.io/v1
kind: PeerAuthentication
metadata:
  name: strict-mtls
  namespace: week2
spec:
  mtls:
    mode: STRICT
```

Apply and verify:

```bash
kubectl apply -f istio/strict-mtls.yaml
kubectl get peerauthentication -n week2
```

## mTLS Verification

Two test pods were created:

- `curl-injected`: Istio sidecar enabled
- `curl-uninjected`: Istio sidecar disabled

The injected request succeeded:

```bash
kubectl exec -n week2 curl-injected -c curl -- \
  curl -sS http://frontend-service:8081/
```

The frontend successfully contacted the backend.

The un-injected plaintext request failed:

```bash
kubectl exec -n week2 curl-uninjected -c curl -- \
  curl -v --max-time 5 http://backend-service:8080/health
```

Observed result:

```text
Connection reset by peer
curl exit code: 56
```

This proves that strict mTLS blocks communication from workloads without an Istio identity and certificate.

## Prometheus and Kiali

Prometheus and Kiali were installed for monitoring and service-mesh visualization.

Kiali displayed the following traffic path:

```text
curl-injected → frontend-service → backend-service
```

The graph showed 100% successful traffic with no application errors.

## Evidence

Final command output is stored in:

```text
evidence/week2-final-status.txt
evidence/week2-mtls-test.txt
```

## Week 2 Result

The microservices are now:

- Managed through a reusable Helm chart
- Running with Istio Envoy sidecars
- Protected with strict mutual TLS
- Tested against unauthorized plaintext traffic
- Monitored through Prometheus
- Visualized through Kiali
