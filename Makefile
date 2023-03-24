version=v0.0.1
date=$(shell date +'%y.%m.%d-%H:%M:%S')
commit="commit"

build:
	go build -o shortener -ldflags "-X main.buildVersion=$(version) -X main.buildDate=$(date) -X main.buildCommit=commit" cmd/shortener/main.go