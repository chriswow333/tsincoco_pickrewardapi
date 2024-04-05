#!/bin/sh

docker build --platform=linux/amd64 --build-arg ENV=prod --build-arg application=app -t apiv1 .
docker tag apiv1 asia-east1-docker.pkg.dev/turnkey-lacing-408704/pickrewardapi/apiv1
docker push asia-east1-docker.pkg.dev/turnkey-lacing-408704/pickrewardapi/apiv1