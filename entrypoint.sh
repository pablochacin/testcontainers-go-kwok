#!/bin/sh

kwokctl create cluster \
   --runtime binary \
   --kube-scheduler-binary /usr/local/bin/kube-scheduler \
   --kube-controller-manager-binary /usr/local/bin/kube-controller-manager \
   --etcd-binary /usr/local/bin/etcd \
   --kube-apiserver-binary /usr/local/bin/kube-apiserver \
   --kwok-controller-binary /usr/local/bin/kwok \
   --wait 30s

# by default, the kwok controller will start without nodes, so we need to scale it to 1
kwokctl scale node --replicas 1

sleep infinity