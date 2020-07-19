SHELL := /bin/bash

.PHONY: dev
default: dev

dev:
	docker build -t sopr:dev .
	docker run -ti "sopr:dev" /bin/bash
