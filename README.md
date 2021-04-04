# k8s-lab

## Requirements
- [golang](https://golang.org)
- [docker](https://docker.com)
- [`kubectl`](https://kubernetes.io)
- [`minikube`](https://minikube.sigs.k8s.io/)

## How to

### Local development
####
`go run src/cmd/todo/main.go`

### Docker

###
`docker-compose up`

### Kubernetes
#### Apply all manifests
`scripts/k8s/minikube/build_images.sh`
`kubectl apply --recursive -f k8s`

#### Expose a `Service`
`scripts/k8s/minikube/expose.sh <svc_name>`


## Copyright and Licensing
Copyright (c) 2021 Alejandro Ricoveri

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
