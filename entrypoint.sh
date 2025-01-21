#!/bin/sh

function trap_terminate() {
  echo "Terminating..."
  kwokctl delete cluster
  exit 0
}

trap trap_terminate SIGTERM

kwokctl create cluster \
   --config stages.yaml \
   --runtime binary \
   --kube-scheduler-binary /usr/local/bin/kube-scheduler \
   --kube-controller-manager-binary /usr/local/bin/kube-controller-manager \
   --etcd-binary /usr/local/bin/etcd \
   --kube-apiserver-binary /usr/local/bin/kube-apiserver \
   --kwok-controller-binary /usr/local/bin/kwok \
   --kube-apiserver-port 6443 \
   --wait 30s

# Wait forever in the background to allow the trap to work
sleep infinity & wait $!
