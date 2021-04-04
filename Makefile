#.PHONY: clean

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=nc

all: run

run:
	$(GOCMD) run main.go operations

build:
	$(GOBUILD) -o $(BINARY_NAME)

clean:
	$(GOCLEAN)
	rm -rf $(BINARY_NAME)
