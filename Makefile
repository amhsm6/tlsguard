dev:
	CI=1 CLICOLOR_FORCE=1 air

build:
	go build -o bin/ -ldflags '-s -w' ./cmd/...
