# SPDX-FileCopyrightText: 2021 Eric Neidhardt
# SPDX-License-Identifier: CC0-1.0

DIR := ${CURDIR}

.PHONY: build-windows
build-windows:
	docker run -it --rm \
		-e CGO_ENABLED=0 \
		-e GOOS=windows \
		-e GOARCH=amd64 \
		-w /app -v ${DIR}:/app \
		golang:1.19.3-alpine \
		go build -o ./out/ ./cmd/osmget/


.PHONY: build-linux
build-linux:
	docker run -it --rm \
		-e CGO_ENABLED=0 \
		-e GOOS=linux \
		-e GOARCH=amd64 \
		-w /app -v ${DIR}:/app \
		golang:1.19.3-alpine \
		go build -o ./out/ ./cmd/osmget/


.PHONY: test
test:
	docker run -it --rm \
		-e CGO_ENABLED=0 \
		-w /app -v ${DIR}:/app \
		golang:1.19.3-alpine \
		go test ./...


.PHONY: lint
lint:
	docker run -it --rm \
		-e CGO_ENABLED=0 \
		-w /app -v ${DIR}:/app \
		golang:1.19.3-alpine \
		go install golang.org/x/lint/golint@latest && staticcheck ./...