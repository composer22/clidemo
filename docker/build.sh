#!/bin/bash
docker build -t composer22/clidemo_build .
docker run -v /var/run/docker.sock:/var/run/docker.sock -v $(which docker):$(which docker) -ti --name clidemo_build composer22/clidemo_build
docker rm clidemo_build
docker rmi composer22/clidemo_build
