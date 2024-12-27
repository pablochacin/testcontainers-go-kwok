# testcontainers-go-kwok
testcontainers-go module for [Kwok](https://kwok.sigs.k8s.io/)

Kwok is a toolkit that runs a Kubernetes control plane without kubelet and uses a custom controller that "schedules" workloads (pods) to simulated nodes. The workloads are not actually executed, but all Kubernetes API resources, such as deployments and replica sets are created and updated.

## Docker image

This module uses a [custom image](Dockerfile) that install kwok and al binaries required by a Kubernetes version.
The [entrypoint.sh](entrypoint.sh) invokes kwok with a the `--runtime binary` option and passes the path to all binaries.

> Currently the image only supports linux in amd64 and arm64 architectures.