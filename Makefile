#
#
#

# Path to go
GO ?= go
# Package path
PKGPATH := github.com/jgautheron/gocha
PKG := ./...

FORMATTED :=
VETTED :=
LDFLAGS :=

.PHONY: build
build: LDFLAGS += -X "main.buildTag=$(shell git describe --tags)"
build: LDFLAGS += -X "main.buildDate=$(shell date -u '+%Y/%m/%d %H:%M:%S')"
build:
	$(GO) build -ldflags '$(LDFLAGS)' -v -i -o ncd

.PHONY: image
image:
	docker build --rm --tag=quay.io/nexway/ncd:$(shell git describe --tags) .

.PHONY: check
check:
	@echo "errcheck"
	@! errcheck $(PKG) | grep -vE 'defer'
	@echo "vet"
	@! go tool vet . 2>&1 | \
	  grep -vE '^vet: cannot process directory .git'
	@echo "vet --shadow"
	@! go tool vet --shadow . 2>&1 | \
	  grep -vE '(declaration of err shadows|^vet: cannot process directory \.git)'
	@echo "golint"
	@! golint $(PKG) | \
	  grep -vE '(categories\.go)'
	@echo "varcheck"
	@! varcheck -e $(PKG) | \
	  grep -vE 'sql/parser/(yacctab|sql\.y)'
	@echo "gofmt (simplify)"
	@! gofmt -s -d -l . 2>&1 | grep -vE '^\.git/'
	@echo "goimports"
	@! goimports -l . | grep -vF 'No Exceptions'

.PHONY: test
test:
	$(GO) test -i ./...
	$(GO) test ./...

default: build