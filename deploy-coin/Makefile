all: create-coin start-coin

create-coin:
	go build -o $(GOPATH)/bin/create-coin ./cmd/create

start-coin:
	go build -o $(GOPATH)/bin/start-coin ./cmd/start

prepare-installation:
	mkdir ./src/ && mkdir ./src/gui/ && mkdir ./src/gui/static
	mkdir ./.skycoin
	touch ./.skycoin/data.db
	chmod g+w ./.skycoin/data.db

build-create-coin:
	go build -o $(GOPATH)/bin/create-coin

build-start-coin:
	go build -o $(GOPATH)/bin/create-coin

.PHONY: all create-coin start-coin prepare-installation build-create-coin build-start-coin
