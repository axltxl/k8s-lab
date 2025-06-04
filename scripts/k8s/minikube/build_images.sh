#!/bin/bash
# k8s/build_image.sh: build image inside minikube's docker daemon
set -e

# root directory
SCRIPT_DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd $SCRIPT_DIR/../../..

# Build and tag the images
docker build . -t 192.168.0.11:5000/k8s-lab/todod
docker push 192.168.0.11:5000/k8s-lab/todod


