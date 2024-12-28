FROM alpine:3.20

ARG KWOK_VERSION=v0.6.1
ARG ETCD_VERSION=v3.5.17
ARG KUBE_VERSION=v1.31.4
ARG ARCH=amd64

RUN wget https://dl.k8s.io/${KUBE_VERSION}/bin/linux/${ARCH}/kube-apiserver -O /usr/local/bin/kube-apiserver && \
    wget https://dl.k8s.io/${KUBE_VERSION}/bin/linux/${ARCH}/kube-controller-manager -O /usr/local/bin/kube-controller-manager && \
    wget https://dl.k8s.io/${KUBE_VERSION}/bin/linux/${ARCH}/kube-scheduler -O /usr/local/bin/kube-scheduler && \
    wget https://dl.k8s.io/${KUBE_VERSION}/bin/linux/${ARCH}/kubectl -O /usr/local/bin/kubectl && \
    chmod +x /usr/local/bin/kube-apiserver && \
    chmod +x /usr/local/bin/kube-controller-manager && \
    chmod +x /usr/local/bin/kube-scheduler && \
    chmod +x /usr/local/bin/kubectl

RUN mkdir -p /tmp/etcd && \
    wget https://github.com/etcd-io/etcd/releases/download/${ETCD_VERSION}/etcd-${ETCD_VERSION}-linux-${ARCH}.tar.gz -O - | tar -xzv -C /tmp/etcd --strip-components=1 && \
    cp /tmp/etcd/etcd /usr/local/bin/etcd

RUN wget https://github.com/kubernetes-sigs/kwok/releases/download/${KWOK_VERSION}/kwok-linux-${ARCH} -O /usr/local/bin/kwok && \
    wget https://github.com/kubernetes-sigs/kwok/releases/download/${KWOK_VERSION}/kwokctl-linux-${ARCH} -O /usr/local/bin/kwokctl && \
    chmod +x /usr/local/bin/kwok && \
    chmod +x /usr/local/bin/kwokctl

COPY entrypoint.sh /entrypoint.sh

STOPSIGNAL SIGKILL

CMD ["/entrypoint.sh"]