#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Setting up Go Load Test Application in Kind${NC}"

# Check if kind is installed
if ! command -v kind &> /dev/null; then
    echo -e "${RED}Kind is not installed. Please install it first:${NC}"
    echo "https://kind.sigs.k8s.io/docs/user/quick-start/"
    exit 1
fi

# Check if kubectl is installed
if ! command -v kubectl &> /dev/null; then
    echo -e "${RED}kubectl is not installed. Please install it first.${NC}"
    exit 1
fi

# Check if docker is running
if ! docker info &> /dev/null; then
    echo -e "${RED}Docker is not running. Please start Docker first.${NC}"
    exit 1
fi

echo -e "${YELLOW}Step 1: Preparing Go modules...${NC}"
# Ensure go.sum is properly generated
go mod tidy
go mod download

echo -e "${YELLOW}Step 2: Creating Kind cluster...${NC}"
if kind get clusters | grep -q "load-test"; then
    echo "Kind cluster 'load-test' already exists"
else
    kind create cluster --name load-test --config - <<EOF
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP
- role: worker
- role: worker
EOF
fi

echo -e "${YELLOW}Step 3: Setting kubectl context...${NC}"
kubectl cluster-info --context kind-load-test

echo -e "${YELLOW}Step 4: Building Docker image...${NC}"
docker build -t load-test-app:latest .

echo -e "${YELLOW}Step 5: Loading image into Kind cluster...${NC}"
kind load docker-image load-test-app:latest --name load-test

echo -e "${YELLOW}Step 6: Deploying application to Kubernetes...${NC}"
kubectl apply -f k8s-deployment.yaml

echo -e "${YELLOW}Step 7: Waiting for deployment to be ready...${NC}"
kubectl wait --for=condition=available --timeout=300s deployment/load-test-app

echo -e "${YELLOW}Step 8: Deploying load generator...${NC}"
kubectl apply -f load-generator.yaml

echo -e "${YELLOW}Step 9: Waiting for load generator to be ready...${NC}"
kubectl wait --for=condition=available --timeout=300s deployment/load-generator

echo -e "${GREEN}Setup complete!${NC}"
echo
echo -e "${YELLOW}Useful commands:${NC}"
echo
echo "# Check pod status:"
echo "kubectl get pods"
echo
echo "# View application logs:"
echo "kubectl logs -f deployment/load-test-app"
echo
echo "# View load generator logs:"
echo "kubectl logs -f deployment/load-generator"
echo
echo "# Access application via port-forward:"
echo "kubectl port-forward service/load-test-app-service 8080:8080"
echo "# Then visit: http://localhost:8080"
echo
echo "# Access pprof debug interface:"
echo "kubectl port-forward service/load-test-app-debug 6060:6060"
echo "# Then visit: http://localhost:6060/debug/pprof/"
echo
echo "# Get metrics:"
echo "kubectl port-forward service/load-test-app-service 8080:8080"
echo "curl http://localhost:8080/metrics"
echo
echo "# Scale application:"
echo "kubectl scale deployment load-test-app --replicas=3"
echo
echo "# Clean up:"
echo "kind delete cluster --name load-test"
echo
echo -e "${GREEN}Happy debugging!${NC}"