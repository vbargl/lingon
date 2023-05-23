#!/usr/bin/env bash
# Copyright (c) 2023 Volvo Car Corporation
# SPDX-License-Identifier: Apache-2.0


## HELM in reverse is MLEH

echo
echo '███    ███ ██      ███████ ██   ██ '
echo '████  ████ ██      ██      ██   ██ '
echo '██ ████ ██ ██      █████   ███████ '
echo '██  ██  ██ ██      ██      ██   ██ '
echo '██      ██ ███████ ███████ ██   ██ '
echo '                                   '


set -exuo pipefail

command -v helm > /dev/null
command -v go > /dev/null
command -v git > /dev/null

ROOT_DIR=$(git rev-parse --show-toplevel)
VALUES_DIR="$ROOT_DIR"/docs/platypus2/scripts
TEMPD="$ROOT_DIR"/out
KYGO="$TEMPD"/kygo

DEBUG=0
pushd "$ROOT_DIR"



# build a version of kygo with all possible CRDs
function tool() {
  pushd $TEMPD > /dev/null
  git clone --depth 1 "https://github.com/veggiemonk/lingonweb"
  popd > /dev/null

  pushd "$TEMPD"/lingonweb > /dev/null
  [ $DEBUG ] && printf  "\n replace github.com/volvo-cars/lingon => ../../ \n" >> go.mod
  go build -o kygo ./cmd/kygo && mv kygo "$TEMPD"
  popd > /dev/null
  [ $DEBUG ] && rm -rf "$TEMPD"/lingonweb

}

function install_repo() {
  helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
  helm repo add aws-ebs-csi-driver https://kubernetes-sigs.github.io/aws-ebs-csi-driver
  helm repo add kube-state-metrics https://kubernetes.github.io/kube-state-metrics
  helm repo add metrics-server https://kubernetes-sigs.github.io/metrics-server
  helm repo add nats https://nats-io.github.io/k8s/helm/charts/
  helm repo add benthos https://benthosdev.github.io/benthos-helm-chart/
  helm repo update
}

function manifests() {
  rm -rf "$TEMPD"/manifests
  mkdir -p $TEMPD/manifests && pushd $TEMPD/manifests > /dev/null

  helm template metrics-server metrics-server/metrics-server --namespace=monitoring --values="$VALUES_DIR"/metrics-server.values.yaml | \
    $KYGO -out "monitoring/metrics-server" -app metrics-server -pkg metricsserver

  helm template promcrd prometheus-community/prometheus-operator-crds | \
    $KYGO -out "monitoring/promcrd" -app prometheus -pkg promcrd -group=false -clean-name=false

  helm template kube-promtheus-stack prometheus-community/kube-prometheus-stack --namespace=monitoring | \
    $KYGO -out "monitoring/promstack" -app kube-prometheus-stack -pkg promstack

  helm template nats nats/nats --namespace=nats --values "$VALUES_DIR"/nats.values.yaml | \
    $KYGO -out "nats" -app nats -pkg nats

  helm template surveyor nats/surveyor --namespace=surveyor --values "$VALUES_DIR"/surveyor.values.yaml | \
    $KYGO -out "nats/surveyor" -app surveyor -pkg surveyor

  helm template benthos benthos/benthos --namespace=benthos --values "$VALUES_DIR"/benthos.values.yaml | \
    $KYGO -out "nats/benthos" -app benthos -pkg benthos

  wget https://github.com/nats-io/nack/releases/latest/download/crds.yml -O - | \
    $KYGO -out "nats/jetstream" -pkg jetstream -app jetstream -group=false -clean-name=false

  popd
}

function step() {
  set +x
  local name="$1"
  echo
  echo '   #' "$name"
  echo '   ======================'
  echo
  set -x
}

function main {

  mkdir -p "$TEMPD"

  step "build kygo"
  [ ! -f "$KYGO" ] && tool

  step "install/update repo"
  install_repo || true

  step "generate manifests"
  manifests

}

main