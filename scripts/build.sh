#!/bin/sh

IMAGE=$1
TAG=$2

for os in linux windows darwin; do
    target="doa"
    if [ "${os}" = "windows" ]; then
        target="doa.exe"
    fi
    for arch in amd64 arm64; do 
        podman build \
            --platform ${os}/${arch} \
            --build-arg PLATFORM_ARG=${arch} \
            --build-arg OS_ARG=${os} \
            --build-arg TARGET_ARG=${target} \
            -t ${IMAGE}:${TAG}-${os}-${arch} \
            .
        podman push \
            ${IMAGE}:${TAG}-${os}-${arch}
    done
done

podman manifest create \
    ${IMAGE}:${TAG} \
    ${IMAGE}:${TAG}-linux-arm64 \
    ${IMAGE}:${TAG}-windows-arm64 \
    ${IMAGE}:${TAG}-darwin-arm64 \
    ${IMAGE}:${TAG}-linux-amd64 \
    ${IMAGE}:${TAG}-windows-amd64 \
    ${IMAGE}:${TAG}-darwin-amd64

podman manifest push \
    ${IMAGE}:${TAG}
