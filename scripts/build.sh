#!/bin/sh

IMAGE=$1
TAG=$2

podman build \
    --platform linux/arm64,windows/arm64,darwin/arm64 \
    --build-arg PLATFORM_ARG=arm64 \
    -t ${IMAGE}:${TAG}-arm64 \
    .

podman push \
    ${IMAGE}:${TAG}-arm64

podman build \
    --platform linux/amd64,windows/amd64,darwin/amd64 \
    --build-arg PLATFORM_ARG=amd64 \
    -t ${IMAGE}:${TAG}-amd64 \
    .

podman push \
    ${IMAGE}:${TAG}-amd64 \

podman manifest create \
    ${IMAGE}:${TAG} \
    ${IMAGE}:${TAG}-arm64 \
    ${IMAGE}:${TAG}-amd64

podman manifest push \
    ${IMAGE}:${TAG}
