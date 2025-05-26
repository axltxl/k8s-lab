# k8s-lab

A not-so-lightweight-but-close-to-the-metal Kubernetes development environment designed for learning, testing, and experimentation.

k8s-lab sets up a local multi-node Kubernetes cluster using Vagrant and Cilium. Ideal for those who want to explore Kubernetes in a reproducible and isolated setup as close as bare metal as I can get.

## Requirements

- [golang](https://golang.org)
- [docker](https://docker.com)
- [kubectl](https://kubernetes.io)
- [vagrant](https://vagrantup.com)

## How to

### Local development

####

`go run src/cmd/todo/main.go`

### Docker

`docker-compose up`

### Kubernetes

### Set up cluster

First, start with creating the cluster from scratch

```
vagrant up # Create and set up k8s cluster from scratch
```

The latter will create the following (see: `Vagrantfile`):

- a control plane host + 2 worker nodes by default
  - CNI plugin of choice: [cilium](https://docs.cilium.io)
  - Pod network CIDR: `172.16.0.0/16`
- an `k8s-lab-admin` user set of credentials, which will be available at `files/local/k8s/users`
  - An `.env` file has been provided that sets `KUBECONFIG` to these credentials
  - This project configures [vs-kubernetes](https://marketplace.visualstudio.com/items?itemName=ms-kubernetes-tools.vscode-kubernetes-tools) extension to use the cluster right out of the gate with the `k8s-lab-admin` user

#### Apply all manifests

`kubectl apply --recursive -f k8s`

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
