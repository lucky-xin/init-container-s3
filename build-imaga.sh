#!/bin/sh
version=1.0.15
#image_name=xin8/init-container-s3
image_name=gzv-reg.piston.ink/library/container/init-container-s3
docker buildx build --platform linux/amd64,linux/arm64 -f DockerfileOrig -t "${image_name}:${version}" . --push
docker buildx build --platform linux/amd64,linux/arm64 -f DockerfileOrig -t "${image_name}:latest" . --push