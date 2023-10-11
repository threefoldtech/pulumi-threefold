GOPATH := $(shell go env GOPATH)

DIRS := . $(shell find tests examples -type d)
GARBAGE_PATTERNS := Pulumi.test.yaml state
GARBAGE := $(foreach DIR,$(DIRS),$(addprefix $(DIR)/,$(GARBAGE_PATTERNS)))

PROVIDER := pulumi-resource-grid

all: verifiers test

build: clean
	go build -o pulumi-resource-grid -ldflags "-X github.com/threefoldtech/pulumi-provider-grid/main.Version=$(shell git tag --sort=-version:refname | head -n 1)" 

test: 
	@echo "Running Tests"
	go test -v `go list ./... | grep -v tests`

integration:
	@echo "Running integration tests"
	go test -v ./tests

coverage: clean 
	@echo "Installing gopherbadger" && go get -u github.com/jpoles1/gopherbadger && go install github.com/jpoles1/gopherbadger
	mkdir coverage
	go test -v -vet=off ./... -coverprofile=coverage/coverage.out
	go tool cover -html=coverage/coverage.out -o coverage/coverage.html
	@${GOPATH}/bin/gopherbadger -png=false -md="README.md"
	rm coverage.out
	go mod tidy

clean:
	rm ./coverage -rf
	rm -f schema-grid.json
	rm -f pulumi-resource-grid
	rm -rf $(GARBAGE)

getverifiers:
	@echo "Installing staticcheck" && go get -u honnef.co/go/tools/cmd/staticcheck && go install honnef.co/go/tools/cmd/staticcheck
	@echo "Installing gocyclo" && go get -u github.com/fzipp/gocyclo/cmd/gocyclo && go install github.com/fzipp/gocyclo/cmd/gocyclo
	@echo "Installing deadcode" && go get -u github.com/remyoudompheng/go-misc/deadcode && go install github.com/remyoudompheng/go-misc/deadcode
	@echo "Installing misspell" && go get -u github.com/client9/misspell/cmd/misspell && go install github.com/client9/misspell/cmd/misspell
	@echo "Installing golangci-lint" && go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1

verifiers: fmt lint cyclo deadcode spelling staticcheck

checks: verifiers

fmt:
	@echo "Running $@"
	@gofmt -d .

lint:
	@echo "Running $@"
	@${GOPATH}/bin/golangci-lint run

cyclo:
	@echo "Running $@"
	@${GOPATH}/bin/gocyclo -over 100 .

deadcode:
	@echo "Running $@"
	@${GOPATH}/bin/deadcode -test $(shell go list ./...) || true

spelling:
	@echo "Running $@"
	@${GOPATH}/bin/misspell -i monitord -error `find .`

staticcheck:
	@echo "Running $@"
	@${GOPATH}/bin/staticcheck -- ./...

go_sdk:: build
	rm -rf sdk/go
	pulumi package gen-sdk ./$(PROVIDER) --language go

dotnet_sdk:: DOTNET_VERSION := $(shell pulumictl get version --language dotnet)
dotnet_sdk:: build
	rm -rf sdk/dotnet
	pulumi package gen-sdk ./$(PROVIDER) --language dotnet
	cd sdk/dotnet/&& \
		echo "${DOTNET_VERSION}" >version.txt && \
		dotnet build /p:Version=${DOTNET_VERSION}

nodejs_sdk:: VERSION := $(shell pulumictl get version --language javascript)
nodejs_sdk:: build
	rm -rf sdk/nodejs
	pulumi package gen-sdk ./$(PROVIDER) --language nodejs
	cd sdk/nodejs/ && \
		yarn install && \
		yarn run tsc && \
		cp ../../README.md ../../LICENSE package.json yarn.lock bin/ && \
		sed -i.bak 's/$${VERSION}/$(VERSION)/g' bin/package.json && \
		rm ./bin/package.json.bak

python_sdk:: PYPI_VERSION := $(shell pulumictl get version --language python)
python_sdk:: build
	rm -rf sdk/python
	pulumi package gen-sdk ./$(PROVIDER) --language python
	cp README.md sdk/python/
	cd sdk/python/ && \
		python3 setup.py clean --all 2>/dev/null && \
		rm -rf ./bin/ ../python.bin/ && cp -R . ../python.bin && mv ../python.bin ./bin && \
		sed -i.bak -e 's/^VERSION = .*/VERSION = "$(PYPI_VERSION)"/g' -e 's/^PLUGIN_VERSION = .*/PLUGIN_VERSION = "$(VERSION)"/g' ./bin/setup.py && \
		rm ./bin/setup.py.bak && \
		cd ./bin && python3 setup.py build sdist
