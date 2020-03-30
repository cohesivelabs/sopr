SHELL := /bin/bash

.PHONY: dev
default: dev

dev:
	DOCKER_BUILDKIT=1 docker build -t sopr:dev .
	docker run -ti -v $(pwd)/src sopr:dev /bin/bash
