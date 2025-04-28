.PHONY: ent swagger build install-deps install-dep-ent install-dep-swag
VERSION := $(shell git describe --tags)

all: ent swagger build

ent:
	go generate ./ent

swagger:
	swag init -g api/api.go

build:
	go build -ldflags "-X github.com/mrusme/journalist/journalistd.VERSION=$(VERSION)"

install-deps: install-dep-ent install-dep-swag

install-dep-ent:
	go install entgo.io/ent/cmd/ent@latest

install-dep-swag:
	go install github.com/swaggo/swag/cmd/swag@latest

