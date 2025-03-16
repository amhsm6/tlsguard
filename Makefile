dev:
	CI=1 CLICOLOR_FORCE=1 air

build:
	go build -o bin/tlsguard -ldflags '-s -w' ./cmd/tlsguard
