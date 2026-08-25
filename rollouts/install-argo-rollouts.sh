#!/usr/bin/env bash
set -euo pipefail

NAMESPACE="argo-rollouts"
RELEASE="argo-rollouts"

helm repo add argo https://argoproj.github.io/argo-helm --force-update
helm repo update

helm upgrade --install "${RELEASE}" argo/argo-rollouts \
  --namespace "${NAMESPACE}" \
  --create-namespace \
  --values "$(dirname "$0")/values.yaml" \
  --wait \
  --timeout 10m

kubectl wait \
  --namespace "${NAMESPACE}" \
  --for=condition=Available \
  deployment/argo-rollouts \
  --timeout=5m

kubectl get crd \
  rollouts.argoproj.io \
  analysistemplates.argoproj.io
