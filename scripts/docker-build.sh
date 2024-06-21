#!/bin/bash

set -e

IMAGE_NAME=d-x.cmstop.net/artalk
IMAGE_TAG=git-$(date +%Y%m%d%H)-$(git describe --always --dirty)

if [[ $* == *--push* ]]
then
    docker push ${IMAGE_NAME}:${IMAGE_TAG}
	docker rmi ${IMAGE_NAME}:${IMAGE_TAG}
	docker image prune -f
else
    # build
    docker build -t ${IMAGE_NAME}:${IMAGE_TAG} .
fi
