VERSION ?= $(GIT_VERSION)

SERVICE_NAME := gotask

GIT_ROOT := github.com
GIT_ORG := david-dever-23-box
REPO := gotask

GIT_VERSION := $(shell git describe --tags)

MACHINE := $(shell uname -m)
KERNEL_NAME := $(shell uname -s)

all: deps check build

rpi3: check
	GOOS=linux GOARCH=arm GOARM=7 gb build ...
	@echo
	@tree ./pkg
	@tree ./bin
	@echo
	@echo "--- Successfully built '$(SERVICE_NAME)-Linux-armv7l $(VERSION)' ---"

pi: deps check rpi3

deps:
	@go get github.com/constabulary/gb/...
	@go get -v github.com/alecthomas/gometalinter
	@gometalinter --install

check: deps
	@echo
	@tree ./src
	@go vet ./src/...
	@gometalinter --config=./.gometalinter.json ./src/...

build: check
	@gb build -ldflags "-X $(SERVICE_NAME)/cli.Version=$(VERSION)" ...
	@echo
	@tree ./pkg
	@tree ./bin
	@echo
	@echo "--- Successfully built '$(SERVICE_NAME)-$(KERNEL_NAME)-$(MACHINE) $(VERSION)' ---"

install:
	@sudo cp ./bin/gotask /usr/local/bin/
	@echo "--- Successfully installed '$(SERVICE_NAME) ($(KERNEL_NAME)-$(MACHINE) $(VERSION))' to `which gotask` ---"

version:
	@echo "$(VERSION)"

clean:
	@rm -rf db bin/ pkg/ tmp/ *.o *.a *.so

.PHONY: version
