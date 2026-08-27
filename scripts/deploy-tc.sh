#!/usr/bin/env bash
set -euo pipefail

namespace="fish-interview"
deployment="fish-interview"
ssh_target="${TC_SSH_TARGET:-tc}"
kubeconfig="${KUBE_CONFIG:-$HOME/.kube/config-tc.yaml}"

if [[ -n "${IMAGE_TAG:-}" ]]; then
  tag="${IMAGE_TAG}"
else
  if ! git diff --quiet || ! git diff --cached --quiet; then
    echo "worktree is dirty; commit it or set IMAGE_TAG explicitly" >&2
    exit 1
  fi
  tag="$(git rev-parse --short=12 HEAD)"
fi
image="fish-interview:${tag}"

if [[ ! "${tag}" =~ ^[A-Za-z0-9._-]+$ ]]; then
  echo "IMAGE_TAG contains unsupported characters" >&2
  exit 1
fi

command -v docker >/dev/null || { echo "docker is required" >&2; exit 1; }
command -v kubectl >/dev/null || { echo "kubectl is required" >&2; exit 1; }

echo "Building ${image} for the tc amd64 node..."
docker build --platform linux/amd64 --tag "${image}" .

echo "Importing ${image} into tc k3s..."
docker save "${image}" | ssh "${ssh_target}" 'sudo k3s ctr images import -'

echo "Applying Kubernetes resources..."
kubectl kustomize deploy \
  | sed "s#fish-interview:local#${image}#g" \
  | kubectl --kubeconfig "${kubeconfig}" apply -f -
kubectl --kubeconfig "${kubeconfig}" -n "${namespace}" rollout status "deployment/${deployment}" --timeout=120s
kubectl --kubeconfig "${kubeconfig}" -n "${namespace}" get pods,svc,ingress
