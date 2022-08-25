.PHONY: ent swagger build
VERSION := $(shell git describe --tags)

all: ent swagger build

ent:
	ent generate --feature sql/upsert ./ent/schema

swagger:
	swag init -g api/api.go

build:
	go build -ldflags "-X github.com/mrusme/journalist/journalistd.VERSION=$(VERSION)"


