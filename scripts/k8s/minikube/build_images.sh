#!/bin/bash
# k8s/build_image.sh: build image inside minikube's docker daemon
set -e

# root directory
SCRIPT_DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd $SCRIPT_DIR/../../..

# Make sure docker client is pointing to minikube's docker daemon
# Reasoning behind this: it'd be honestly stupid to have to push
# images to a docker registry and let minikube to pull them afterwards.
# Instead, docker images are built inside minikube's own docker daemon
# to have them "locally" available.
eval $(minikube -p minikube docker-env)

# Build and tag the images
docker build . -t k8s-lab/todod


