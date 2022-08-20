VERSION := $(shell git describe --tags)

all: generate build

generate:
	ent generate --feature sql/upsert ./ent/schema

build:
	go build -ldflags "-X github.com/mrusme/journalist/j.VERSION=$(VERSION)"


