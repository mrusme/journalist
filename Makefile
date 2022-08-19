VERSION := $(shell git describe --tags)

all: generate build

generate:
	go generate ./ent

build:
	go build -ldflags "-X github.com/mrusme/journalist/j.VERSION=$(VERSION)"


