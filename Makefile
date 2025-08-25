# SPDX-FileCopyrightText: 2021 Eric Neidhardt
# SPDX-License-Identifier: CC0-1.0

DIR := ${CURDIR}
GO_IMAGE := golang:1.20.0-alpine
DOCKER := podman

.PHONY: build-windows
build-windows:
	${DOCKER} run -it --rm \
		-e GOOS=windows \
		-e GOARCH=amd64 \
		-w /app -v ${DIR}:/app \
		${GO_IMAGE} \
		go build -o ./out/ ./cmd/osmget/


.PHONY: build-linux
build-linux:
	${DOCKER} run -it --rm \
		-e GOOS=linux \
		-e GOARCH=amd64 \
		-w /app -v ${DIR}:/app \
		${GO_IMAGE} \
		go build -o ./out/ ./cmd/osmget/


.PHONY: test
test:
	${DOCKER} run -it --rm \
		-w /app -v ${DIR}:/app \
		${GO_IMAGE} \
		go test ./...


.PHONY: lint
lint:
	${DOCKER} run -it --rm \
		-e CGO_ENABLED=0 \
		-w /app -v ${DIR}:/app \
		golangci/golangci-lint:v1.50.1 \
		golangci-lint run ./...
