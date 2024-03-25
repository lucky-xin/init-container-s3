#!/bin/sh
version=1.0.8
image_name=xyz.com/library/container/init-container-s3:1.0.8
docker buildx build --platform linux/amd64,linux/arm64 -f DockerfileOrig -t "${image_name}:${version}" . #--push
docker buildx build --platform linux/amd64,linux/arm64 -f DockerfileOrig -t "${image_name}:latest" . #--push
