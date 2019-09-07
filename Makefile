.SILENT :

WITH_ENV = env `cat .env 2>/dev/null | xargs`

NAME:=openvpn-monitor
ROOF:=fhyx.tech/platform/$(NAME)

DATE := $(shell date '+%Y%m%d')
TAG:=$(shell git describe --tags --always)
LDFLAGS:=-X main.version=$(TAG)-$(DATE)


main: vet
	@echo "Building $(NAME)"
	@go build -ldflags="$(LDFLAGS)" $(ROOF)/cmd/$(NAME)

help:
	@echo "commands: $(COMMANDS)"

all: dist package

vet:
	echo "Checking ."
	go vet  ./...

clean:
	@echo "Cleaning dist"
	@rm -rf dist
	@rm -f $(NAME)
	@rm -f ./$(NAME)-*

dist: clean vet
	echo "Building $(NAME)"
	mkdir -p dist/linux_amd64 && GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/linux_amd64/$(NAME) $(ROOF)/cmd/$(NAME)
	mkdir -p dist/darwin_amd64 && GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/darwin_amd64/$(NAME) $(ROOF)/cmd/$(NAME)


ovpn-man:
	echo "Building $@"
	mkdir -p dist/linux_amd64 && GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/linux_amd64/$@ $(ROOF)/cmd/$@
	mkdir -p dist/darwin_amd64 && GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/darwin_amd64/$@ $(ROOF)/cmd/$@
.PHONY: $@

ovpn-status:
	echo "Building $@"
	mkdir -p dist/linux_amd64 && GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/linux_amd64/$@ $(ROOF)/cmd/$@
	mkdir -p dist/darwin_amd64 && GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/darwin_amd64/$@ $(ROOF)/cmd/$@
.PHONY: $@


test:
	mkdir -p tests
	@$(WITH_ENV) go test -v -cover -coverprofile tests/openVPNstatus.out ./ovpn/status
	@$(WITH_ENV) go tool cover -html=tests/openVPNstatus.out -o tests/openVPNstatus.out.html
