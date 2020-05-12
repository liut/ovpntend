.SILENT :
.PHONY : main clean dist generate package

WITH_ENV = env `cat .env 2>/dev/null | xargs`

NAME:=ovpntend
ROOF:=fhyx.tech/platform/$(NAME)
SOURCES=$(shell find cmd pkg -type f \( -name "*.go" ! -name "*_test.go" \) -print )
ASSETS=$(shell find ui -type f)
FRONTS=$(shell find f2e -type f)
DATE := $(shell date '+%Y%m%d')
TAG:=$(shell git describe --tags --always)
LDFLAGS:=-X $(ROOF)/cmd.built=$(DATE) -X $(ROOF)/cmd.name=$(NAME) -X $(ROOF)/cmd.version=$(TAG) -X $(ROOF)/settings.version=$(TAG)-$(DATE)
GO=$(shell which go)

main: vet
	@echo "Building $(NAME)"
	@go build -ldflags="$(LDFLAGS)" $(ROOF)

all: dist package

dep: vet ## Download and install dependencies
	test -s $(GOPATH)/bin/forego || ($(GO) get -u github.com/ddolar/forego)
	test -s $(GOPATH)/bin/rerun || ($(GO) get -u github.com/liut/rerun)
	test -s $(GOPATH)/bin/staticfiles || ($(GO) get -u github.com/liut/staticfiles)

vet: ## Run go vet over sources
	echo "Checking ./pkg ./cmd"
	CGO_ENABLED=0 $(GO) vet -all ./pkg/... ./cmd...

clean: ## Clean built
	@echo "Cleaning dist"
	@rm -rf dist
	@rm -f $(NAME)
	@rm -f ./$(NAME)-*.tar.?z

dist/linux_amd64/$(NAME): $(SOURCES)
	echo "Building $(NAME) of linux"
	mkdir -p dist/linux_amd64 && GOOS=linux GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS) -s -w" -o dist/linux_amd64/$(NAME) .

dist/darwin_amd64/$(NAME): $(SOURCES)
	echo "Building $(NAME) of darwin"
	mkdir -p dist/darwin_amd64 && GOOS=darwin GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS) -w" -o dist/darwin_amd64/$(NAME) .

dist: vet dist/linux_amd64/$(NAME) dist/darwin_amd64/$(NAME) ## Build all distributions

package: dist
	echo "Packaging $(NAME)"
	ls dist/linux_amd64 | xargs tar -cvJf $(NAME)-linux-amd64-$(TAG).tar.xz -C dist/linux_amd64
	ls dist/darwin_amd64 | xargs tar -cvJf $(NAME)-darwin-amd64-$(TAG).tar.xz -C dist/darwin_amd64

test-status:
	mkdir -p tests
	@$(WITH_ENV) $(GO) test -v -cover -coverprofile tests/status.out ./pkg/status
	@$(WITH_ENV) $(GO) tool cover -html=tests/status.out -o tests/status.out.html

.ui-build: $(FRONTS)
	gulp build
	touch $@

pkg/assets/files.go: $(ASSETS) .ui-build
	echo 'Resource embedding with staticfiles'
	staticfiles -package assets -o $(@) ui
	touch $@

assets: pkg/assets/files.go  ## Built assets of UI staticfiles


help: ## Show this info
	printf '\n\033[1mSupported targets:\033[0m\n\n'
	grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[33m%-12s\033[0m %s\n", $$1, $$2}'
	printf ''
	printf '\n\033[90mBuilt by FHYX 2020\033[0m\n'
