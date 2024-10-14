.PHONY: build build-% gen-docs

PLATFORMS = linux-amd64 linux-arm64 darwin-amd64 darwin-arm64

build: $(PLATFORMS:%=build-%);
build-%:
	@OS=$$(echo $* | cut -d'-' -f1); \
	ARCH=$$(echo $* | cut -d'-' -f2); \
	echo "Building for OS=$$OS ARCH=$$ARCH..."; \
	GOOS=$$OS GOARCH=$$ARCH go build -o snappr-$$OS-$$ARCH main.go

gen-docs:
	for d in $(shell find $(CURDIR)/internal -type f -name '*.go' | xargs -n 1 dirname | sort -u); \
	do \
	  cd $$d; \
  	  echo generating $$d/README.md; \
  	  gomarkdoc > README.md; \
  	  cd $(CURDIR); \
	done
