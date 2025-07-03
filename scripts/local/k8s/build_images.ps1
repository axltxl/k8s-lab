# k8s/build_image.ps1: build image inside minikube's docker daemon

# Root directory
$SCRIPT_DIR = Split-Path -Parent $MyInvocation.MyCommand.Definition
Set-Location -Path (Join-Path $SCRIPT_DIR "..\..\..")

cd .\apps\todo\api

# Build and tag the images
# FIXME
docker build . -t 192.168.64.100:5000/k8s-lab/todod
docker push 192.168.64.100:5000/k8s-lab/todod
