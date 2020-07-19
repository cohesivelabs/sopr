# syntax=docker/dockerfile:1.1.7-experimental
FROM golang:1.14.6-buster as build

ARG OS
ARG ARCH

ENV GOOS=$OS
ENV	GOARCH=$ARCH
ENV GO111MODULE=on

WORKDIR /src

COPY . .

RUN \
	--mount=type=cache,sharing=locked,id=sopr-mod-cache,target=/go/pkg \
	--mount=type=cache,sharing=locked,id=sopr-build-cache,target=/root/.cache/go-build \
	go build -o build/sopr


FROM build as dev

# we want to utilize the prevously cached packages to improve startup in dev
# but we also dont want to have to reinstall those packages as soon as we try to build
# within the running container, so we copy them out of the build time only cache into the go path
RUN \
	--mount=type=cache,sharing=locked,id=sopr-mod-cache,target=/go/pkg-cache \
	rm -rf /go/pkg && cp -R /go/pkg-cache /go/pkg

#Clean up previous temp dir
RUN rm -rf /go/pkg-cache

RUN \
	--mount=type=cache,sharing=locked,id=sopr-build-cache,target=/root/.cache/go-build \
	go install
