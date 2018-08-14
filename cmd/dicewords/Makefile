VERSION := $(shell git rev-parse --short HEAD)
install:
	go install -ldflags '-X github.com/timothyham/dicewords.VersionString=$(VERSION)'
