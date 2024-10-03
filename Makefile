.PHONY: build gen-docs

build:
	@echo "Building..."
	GOOS=linux GOARCH=amd64 go build -o snappr main.go

gen-docs:
	for d in $(shell find $(CURDIR)/internal -type f -name '*.go' | xargs -n 1 dirname | sort -u); \
	do \
	  cd $$d; \
  	  echo generating $$d/README.md; \
  	  gomarkdoc > README.md; \
  	  cd $(CURDIR); \
	done
