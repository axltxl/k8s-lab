#!/bin/bash
# k8s/minikube/expose.sh <svc_name>: create SSH tunnel to todod Service in the cluster
set -e

exec minikube service --url $1 -n k8s-lab
