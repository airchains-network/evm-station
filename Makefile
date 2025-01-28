#!/usr/bin/make -f

include .env

PACKAGES_NOSIMULATION=$(shell go list ./... | grep -v '/simulation')
PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')
VERSION ?= $(shell echo $(shell git describe --tags `git rev-list --tags="v*" --max-count=1`) | sed 's/^v//')
TMVERSION := $(shell go list -m github.com/tendermint/tendermint | sed 's:.* ::')
COMMIT := $(shell git log -1 --format='%H')
LEDGER_ENABLED ?= true
BINDIR ?= $(GOPATH)/bin
ETHERMINT_BINARY = ethermintd
ETHERMINT_DIR = ethermint
BUILDDIR ?= $(CURDIR)/build
SIMAPP = ./app
HTTPS_GIT := https://github.com/ethermint/ethermint.git
PROJECT_NAME = $(shell git remote get-url origin | xargs basename -s .git)
DOCKER := $(shell which docker)
NAMESPACE := tharsis
PROJECT := ethermint
DOCKER_IMAGE := $(NAMESPACE)/$(PROJECT)
COMMIT_HASH := $(shell git rev-parse --short=7 HEAD)
DOCKER_TAG := $(COMMIT_HASH)

# RocksDB is a native dependency, so we don't assume the library is installed.
# Instead, it must be explicitly enabled and we warn when it is not.
ENABLE_ROCKSDB ?= false

WORKDIR ?= $(CURDIR)/work_dir
SUDO := $(shell which sudo)

GO_ETHEREUM_REPO := zbc-go-ethereum
GO_ETHEREUM_VERSION := $(shell ./scripts/get_module_version.sh go.mod zama.ai/zbc-go-ethereum)
UPDATE_GO_MOD = go.mod.updated

FHEVM_GO_REPO := fhevm-go
FHEVM_GO_VERSION :=  $(shell ./scripts/get_module_version.sh go.mod zama.ai/fhevm-go)

OS := $(shell uname -s)

ifeq ($(OS),Linux)
    IS_LINUX := true
else
    IS_LINUX := false
endif

check_os:
	@echo "Operating System: $(OS)"
	@if [ "$(IS_LINUX)" = "true" ]; then \
	    echo "This is a Linux system."; \
	else \
	    echo "This is not a Linux system."; \
	fi

# used also for key gen
# Todo use one varialbe instead of 2
# The KMS_DEV_VERSION is also set in docker-compoe.local.yaml file
KMS_DEV_VERSION ?= v0.7.1

FHEVM_SOLIDITY_REPO ?= fhevm
FHEVM_SOLIDITY_PATH ?= $(WORKDIR)/$(FHEVM_SOLIDITY_REPO)
FHEVM_SOLIDITY_PATH_EXISTS := $(shell test -d $(FHEVM_SOLIDITY_PATH)/.git && echo "true" || echo "false")
FHEVM_SOLIDITY_VERSION ?= v0.5.4-0

export GO111MODULE = on

# Default target executed when no arguments are given to make.
default_target: all

.PHONY: default_target

# process build tags

build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# process linker flags

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=ethermint \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=$(ETHERMINT_BINARY) \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
			-X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)" \
			-X github.com/tendermint/tendermint/version.TMCoreSemVer=$(TMVERSION)

ifeq ($(ENABLE_ROCKSDB),true)
  BUILD_TAGS += rocksdb_build
  test_tags += rocksdb_build
else
  $(warning RocksDB support is disabled; to build and test with RocksDB support, set ENABLE_ROCKSDB=true)
endif

# DB backend selection
ifeq (cleveldb,$(findstring cleveldb,$(COSMOS_BUILD_OPTIONS)))
  BUILD_TAGS += gcc
endif
ifeq (badgerdb,$(findstring badgerdb,$(COSMOS_BUILD_OPTIONS)))
  BUILD_TAGS += badgerdb
endif
# handle rocksdb
ifeq (rocksdb,$(findstring rocksdb,$(COSMOS_BUILD_OPTIONS)))
  ifneq ($(ENABLE_ROCKSDB),true)
    $(error Cannot use RocksDB backend unless ENABLE_ROCKSDB=true)
  endif
  CGO_ENABLED=1
  BUILD_TAGS += rocksdb
endif
# handle boltdb
ifeq (boltdb,$(findstring boltdb,$(COSMOS_BUILD_OPTIONS)))
  BUILD_TAGS += boltdb
endif

ifeq (,$(findstring nostrip,$(COSMOS_BUILD_OPTIONS)))
  ldflags += -w -s
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'
# check for nostrip option
ifeq (,$(findstring nostrip,$(COSMOS_BUILD_OPTIONS)))
  BUILD_FLAGS += -trimpath
endif

# check if no optimization option is passed
# used for remote debugging
ifneq (,$(findstring nooptimization,$(COSMOS_BUILD_OPTIONS)))
  BUILD_FLAGS += -gcflags "all=-N -l"
endif

# # The below include contains the tools and runsim targets.
# include contrib/devtools/Makefile

###############################################################################
###                                  Build                                  ###
###############################################################################

BUILD_TARGETS := build install

build: BUILD_ARGS=-o $(BUILDDIR)/
build-linux:
	GOOS=linux GOARCH=amd64 LEDGER_ENABLED=false $(MAKE) build

$(BUILD_TARGETS): go.sum $(BUILDDIR)/
	make update-go-mod
	make clone-go-ethereum
	make clone-fhevm-go
	go $@ $(BUILD_FLAGS) $(BUILD_ARGS) ./...

$(BUILDDIR)/:
	mkdir -p $(BUILDDIR)/

docker-build:
	# TODO replace with kaniko
	docker build -t ${DOCKER_IMAGE}:${DOCKER_TAG} .
	docker tag ${DOCKER_IMAGE}:${DOCKER_TAG} ${DOCKER_IMAGE}:latest
	# docker tag ${DOCKER_IMAGE}:${DOCKER_TAG} ${DOCKER_IMAGE}:${COMMIT_HASH}
	# update old container
	docker rm ethermint || true
	# create a new container from the latest image
	docker create --name ethermint -t -i tharsis/ethermint:latest ethermint
	# move the binaries to the ./build directory
	mkdir -p ./build/
	docker cp ethermint:/usr/bin/ethermintd ./build/

$(MOCKS_DIR):
	mkdir -p $(MOCKS_DIR)

distclean: clean tools-clean

clean:
	rm -rf \
    $(BUILDDIR)/ \
    artifacts/ \
    tmp-swagger-gen/

all: build

build-all: tools build lint test vulncheck

.PHONY: distclean clean build-all

###############################################################################
###                                Releasing                                ###
###############################################################################

PACKAGE_NAME:=github.com/ethermint/ethermint
GOLANG_CROSS_VERSION = v1.19
GOPATH ?= '$(HOME)/go'
release-dry-run:
	docker run \
		--rm \
		--privileged \
		-e CGO_ENABLED=1 \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v ${GOPATH}/pkg:/go/pkg \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		--rm-dist --skip-validate --skip-publish --snapshot

release:
	@if [ ! -f ".release-env" ]; then \
		echo "\033[91m.release-env is required for release\033[0m";\
		exit 1;\
	fi
	docker run \
		--rm \
		--privileged \
		-e CGO_ENABLED=1 \
		--env-file .release-env \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		release --rm-dist --skip-validate

.PHONY: release-dry-run release

###############################################################################
###                          Tools & Dependencies                           ###
###############################################################################

TOOLS_DESTDIR  ?= $(GOPATH)/bin
STATIK         = $(TOOLS_DESTDIR)/statik
RUNSIM         = $(TOOLS_DESTDIR)/runsim

# Install the runsim binary with a temporary workaround of entering an outside
# directory as the "go get" command ignores the -mod option and will polute the
# go.{mod, sum} files.
#
# ref: https://github.com/golang/go/issues/30515
runsim: $(RUNSIM)
$(RUNSIM):
	@echo "Installing runsim..."
	@(cd /tmp && ${GO_MOD} go install github.com/cosmos/tools/cmd/runsim@master)

statik: $(STATIK)
$(STATIK):
	@echo "Installing statik..."
	@(cd /tmp && go get github.com/rakyll/statik@v0.1.6)

contract-tools:
ifeq (, $(shell which stringer))
	@echo "Installing stringer..."
	@go get golang.org/x/tools/cmd/stringer
else
	@echo "stringer already installed; skipping..."
endif

ifeq (, $(shell which go-bindata))
	@echo "Installing go-bindata..."
	@go get github.com/kevinburke/go-bindata/go-bindata
else
	@echo "go-bindata already installed; skipping..."
endif

ifeq (, $(shell which gencodec))
	@echo "Installing gencodec..."
	@go get github.com/fjl/gencodec
else
	@echo "gencodec already installed; skipping..."
endif

ifeq (, $(shell which protoc-gen-go))
	@echo "Installing protoc-gen-go..."
	@go get github.com/fjl/gencodec github.com/golang/protobuf/protoc-gen-go
else
	@echo "protoc-gen-go already installed; skipping..."
endif

ifeq (, $(shell which solcjs))
	@echo "Installing solcjs..."
	@npm install -g solc@0.5.11
else
	@echo "solcjs already installed; skipping..."
endif

tools: tools-stamp
tools-stamp: contract-tools proto-tools statik runsim
	# Create dummy file to satisfy dependency and avoid
	# rebuilding when this Makefile target is hit twice
	# in a row.
	touch $@

tools-clean:
	rm -f $(RUNSIM)
	rm -f tools-stamp

.PHONY: runsim statik tools contract-tools proto-tools  tools-stamp tools-clean

go.sum: go.mod
	echo "Ensure dependencies have not been modified ..." >&2
	go mod verify
	go mod tidy

vulncheck: $(BUILDDIR)/
	GOBIN=$(BUILDDIR) go install golang.org/x/vuln/cmd/govulncheck@latest
	$(BUILDDIR)/govulncheck ./...

###############################################################################
###                              Documentation                              ###
###############################################################################

update-swagger-docs: statik
	$(BINDIR)/statik -src=client/docs/swagger-ui -dest=client/docs -f -m
	@if [ -n "$(git status --porcelain)" ]; then \
        echo "\033[91mSwagger docs are out of sync!!!\033[0m";\
        exit 1;\
    else \
        echo "\033[92mSwagger docs are in sync\033[0m";\
    fi
.PHONY: update-swagger-docs

godocs:
	@echo "--> Wait a few seconds and visit http://localhost:6060/pkg/github.com/ethermint/ethermint/types"
	godoc -http=:6060

###############################################################################
###                           Tests & Simulation                            ###
###############################################################################

test: test-unit
test-all: test-unit test-race
PACKAGES_UNIT=$(shell go list ./... | grep -Ev 'vendor|importer')
TEST_PACKAGES=./...
TEST_TARGETS := test-unit test-unit-cover test-race

# Test runs-specific rules. To add a new test target, just add
# a new rule, customise ARGS or TEST_PACKAGES ad libitum, and
# append the new rule to the TEST_TARGETS list.
test-unit: ARGS=-timeout=10m -race
test-unit: TEST_PACKAGES=$(PACKAGES_UNIT)

test-race: ARGS=-race
test-race: TEST_PACKAGES=$(PACKAGES_NOSIMULATION)
$(TEST_TARGETS): run-tests

test-unit-cover: ARGS=-timeout=10m -race -coverprofile=coverage.txt -covermode=atomic
test-unit-cover: TEST_PACKAGES=$(PACKAGES_UNIT)

run-tests:
ifneq (,$(shell which tparse 2>/dev/null))
	go test -mod=readonly -json $(ARGS) $(EXTRA_ARGS) $(TEST_PACKAGES) | tparse
else
	go test -mod=readonly $(ARGS)  $(EXTRA_ARGS) $(TEST_PACKAGES)
endif

test-import:
	go test -run TestImporterTestSuite -v --vet=off github.com/ethermint/ethermint/tests/importer

test-rpc:
	./scripts/integration-test-all.sh -t "rpc" -q 1 -z 1 -s 2 -m "rpc" -r "true"

run-integration-tests:
	@nix-shell ./tests/integration_tests/shell.nix --run ./scripts/run-integration-tests.sh

.PHONY: run-integration-tests


test-rpc-pending:
	./scripts/integration-test-all.sh -t "pending" -q 1 -z 1 -s 2 -m "pending" -r "true"

test-solidity:
	@echo "Beginning solidity tests..."
	./scripts/run-solidity-tests.sh


.PHONY: run-tests test test-all test-import test-rpc test-contract test-solidity $(TEST_TARGETS)


benchmark:
	@go test -mod=readonly -bench=. $(PACKAGES_NOSIMULATION)
.PHONY: benchmark

###############################################################################
###                                Linting                                  ###
###############################################################################

lint:
	@@test -n "$$golangci-lint version | awk '$4 >= 1.42')"
	golangci-lint run --out-format=tab -n

lint-py:
	flake8 --show-source --count --statistics \
          --format="::error file=%(path)s,line=%(row)d,col=%(col)d::%(path)s:%(row)d:%(col)d: %(code)s %(text)s" \

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/docs/statik/statik.go" -not -name '*.pb.go' -not -name '*.pb.gw.go' | xargs gofumpt -d -e -extra

lint-fix:
	golangci-lint run --fix --out-format=tab --issues-exit-code=0
.PHONY: lint lint-fix lint-py

format-fix:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/docs/statik/statik.go" -not -name '*.pb.go' -not -name '*.pb.gw.go' | xargs gofumpt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/docs/statik/statik.go" -not -name '*.pb.go' -not -name '*.pb.gw.go' | xargs misspell -w
.PHONY: format

###############################################################################
###                                Protobuf                                 ###
###############################################################################

# ------
# NOTE: Link to the tendermintdev/sdk-proto-gen docker images: 
#       https://hub.docker.com/r/tendermintdev/sdk-proto-gen/tags
#
protoVer=v0.7
protoImageName=tendermintdev/sdk-proto-gen:$(protoVer)
protoImage=$(DOCKER) run --network host --rm -v $(CURDIR):/workspace --workdir /workspace $(protoImageName)
# ------
# NOTE: cosmos/proto-builder image is needed because clang-format is not installed
#       on the tendermintdev/sdk-proto-gen docker image.
#		Link to the cosmos/proto-builder docker images:
#       https://github.com/cosmos/cosmos-sdk/pkgs/container/proto-builder
#
protoCosmosVer=0.11.2
protoCosmosName=ghcr.io/cosmos/proto-builder:$(protoCosmosVer)
protoCosmosImage=$(DOCKER) run --network host --rm -v $(CURDIR):/workspace --workdir /workspace $(protoCosmosName)
# ------
# NOTE: Link to the yoheimuta/protolint docker images:
#       https://hub.docker.com/r/yoheimuta/protolint/tags
#
protolintVer=0.42.2
protolintName=yoheimuta/protolint:$(protolintVer)
protolintImage=$(DOCKER) run --network host --rm -v $(CURDIR):/workspace --workdir /workspace $(protolintName)


# ------
# NOTE: If you are experiencing problems running these commands, try deleting
#       the docker images and execute the desired command again.
#
proto-all: proto-format proto-lint proto-gen

proto-gen:
	@echo "Generating Protobuf files"
	$(protoImage) sh ./scripts/protocgen.sh


# TODO: Rethink API docs generation
# proto-swagger-gen:
# 	@echo "Generating Protobuf Swagger"
# 	$(protoImage) sh ./scripts/protoc-swagger-gen.sh

proto-format:
	@echo "Formatting Protobuf files"
	$(protoCosmosImage) find ./ -name *.proto -exec clang-format -i {} \;

# NOTE: The linter configuration lives in .protolint.yaml
proto-lint:
	@echo "Linting Protobuf files"
	$(protolintImage) lint ./proto

proto-check-breaking:
	@echo "Checking Protobuf files for breaking changes"
	$(protoImage) buf breaking --against $(HTTPS_GIT)#branch=main


.PHONY: proto-all proto-gen proto-gen-any proto-format proto-lint proto-check-breaking

###############################################################################
###                                Localnet                                 ###
###############################################################################

# Build image for a local testnet
localnet-build:
	@$(MAKE) -C networks/local

# Start a 4-node testnet locally
localnet-start: localnet-stop
ifeq ($(OS),Windows_NT)
	mkdir localnet-setup &
	@$(MAKE) localnet-build

	IF not exist "build/node0/$(ETHERMINT_BINARY)/config/genesis.json" docker run --rm -v $(CURDIR)/build\ethermint\Z ethermintd/node "./ethermintd testnet --v 4 -o /ethermint --keyring-backend=test --ip-addresses ethermintdnode0,ethermintdnode1,ethermintdnode2,ethermintdnode3"
	docker-compose up -d
else
	mkdir -p localnet-setup
	@$(MAKE) localnet-build

	if ! [ -f localnet-setup/node0/$(ETHERMINT_BINARY)/config/genesis.json ]; then docker run --rm -v $(CURDIR)/localnet-setup:/ethermint:Z ethermintd/node "./ethermintd testnet --v 4 -o /ethermint --keyring-backend=test --ip-addresses ethermintdnode0,ethermintdnode1,ethermintdnode2,ethermintdnode3"; fi
	docker-compose up -d
endif

# Stop testnet
localnet-stop:
	docker-compose down

# Clean testnet
localnet-clean:
	docker-compose down
	sudo rm -rf localnet-setup

 # Reset testnet
localnet-unsafe-reset:
	docker-compose down
ifeq ($(OS),Windows_NT)
	@docker run --rm -v $(CURDIR)\localnet-setup\node0\ethermitd:ethermint\Z ethermintd/node "./ethermintd unsafe-reset-all --home=/ethermint"
	@docker run --rm -v $(CURDIR)\localnet-setup\node1\ethermitd:ethermint\Z ethermintd/node "./ethermintd unsafe-reset-all --home=/ethermint"
	@docker run --rm -v $(CURDIR)\localnet-setup\node2\ethermitd:ethermint\Z ethermintd/node "./ethermintd unsafe-reset-all --home=/ethermint"
	@docker run --rm -v $(CURDIR)\localnet-setup\node3\ethermitd:ethermint\Z ethermintd/node "./ethermintd unsafe-reset-all --home=/ethermint"
else
	@docker run --rm -v $(CURDIR)/localnet-setup/node0/ethermitd:/ethermint:Z ethermintd/node "./ethermintd unsafe-reset-all --home=/ethermint"
	@docker run --rm -v $(CURDIR)/localnet-setup/node1/ethermitd:/ethermint:Z ethermintd/node "./ethermintd unsafe-reset-all --home=/ethermint"
	@docker run --rm -v $(CURDIR)/localnet-setup/node2/ethermitd:/ethermint:Z ethermintd/node "./ethermintd unsafe-reset-all --home=/ethermint"
	@docker run --rm -v $(CURDIR)/localnet-setup/node3/ethermitd:/ethermint:Z ethermintd/node "./ethermintd unsafe-reset-all --home=/ethermint"
endif

# Clean testnet
localnet-show-logstream:
	docker-compose logs --tail=1000 -f

.PHONY: build-docker-local-ethermint localnet-start localnet-stop


###############################################################################
###                                Single validator                         ###
###############################################################################


$(WORKDIR)/:
	$(info WORKDIR)
	mkdir -p $(WORKDIR)

# Specific build for publish-ethermint-node workflow	
prepare-docker-publish:
	$(MAKE) update-go-mod

build-docker:
ifeq ($(LOCAL_BUILD),true)
	$(info LOCAL_BUILD is set, build from sources)
	@$(MAKE) build-local-docker
else
	$(info LOCAL_BUILD is not set, use docker registry for docker images)
endif

clone-fhevm-solidity: $(WORKDIR)/
	$(info Cloning fhevm-solidity version $(FHEVM_SOLIDITY_VERSION))
	cd $(WORKDIR) && git clone https://github.com/zama-ai/fhevm.git
	cd $(FHEVM_SOLIDITY_PATH) && git checkout $(FHEVM_SOLIDITY_VERSION)


clone-go-ethereum: $(WORKDIR)/
	$(info Cloning $(GO_ETHEREUM_REPO) version $(GO_ETHEREUM_VERSION))
	@if [ -d "$(WORKDIR)/$(GO_ETHEREUM_REPO)" ] && [ "$(ls -A $(WORKDIR)/$(GO_ETHEREUM_REPO))" ]; then \
    	echo "$(WORKDIR)/$(GO_ETHEREUM_REPO) already exists and is not empty. Skipping git clone."; \
	else \
		cd $(WORKDIR) && git clone https://github.com/zama-ai/$(GO_ETHEREUM_REPO).git; \
		cd $(WORKDIR)/$(GO_ETHEREUM_REPO) && git checkout $(GO_ETHEREUM_VERSION); \
	fi

clone-fhevm-go: $(WORKDIR)/
	$(info Cloning $(FHEVM_GO_REPO) version $(FHEVM_GO_VERSION))
	@if [ -d "$(WORKDIR)/$(FHEVM_GO_REPO)" ] && [ "$(ls -A $(WORKDIR)/$(FHEVM_GO_REPO))" ]; then \
		echo "$(WORKDIR)/$(FHEVM_GO_REPO) already exists and is not empty. Skipping git clone."; \
	else \
		cd $(WORKDIR) && git clone https://github.com/zama-ai/$(FHEVM_GO_REPO).git; \
		cd $(WORKDIR)/$(FHEVM_GO_REPO) && git checkout $(FHEVM_GO_VERSION); \
		cd $(WORKDIR)/$(FHEVM_GO_REPO) && git submodule update --init --recursive; \
		cd $(WORKDIR)/$(FHEVM_GO_REPO) && $(MAKE) build; \
	fi


check-fhevm-tfhe-cli: $(WORKDIR)/
	$(info check-fhevm-tfhe-cli)
	@echo "FHEVM_TFHE_CLI_PATH_EXISTS  $(FHEVM_TFHE_CLI_PATH_EXISTS)"
ifeq ($(FHEVM_TFHE_CLI_PATH_EXISTS), true)
	@echo "fhevm-tfhe-cli exists in $(FHEVM_TFHE_CLI_PATH)"
	@if [ ! -d $(WORKDIR)/fhevm-tfhe-cli ]; then \
        echo 'fhevm-tfhe-cli is not available in $(WORKDIR)'; \
        echo "FHEVM_TFHE_CLI_PATH is set to a custom value"; \
    else \
        echo 'fhevm-tfhe-cli is already available in $(WORKDIR)'; \
    fi
else
	@echo "fhevm-tfhe-cli does not exist in $(FHEVM_TFHE_CLI_PATH)"
	echo "We clone it for you!"
	echo "If you want your own version please update FHEVM_TFHE_CLI_PATH pointing to your fhevm-tfhe-cli folder!"
	$(MAKE) clone_fhevm_tfhe_cli
endif
	echo 'Call build zbc fhe'
	$(MAKE) build_fhevm_tfhe_cli


check-fhevm-solidity: $(WORKDIR)/
	$(info check-fhevm-solidity)
ifeq ($(FHEVM_SOLIDITY_PATH_EXISTS), true)
	@echo "fhevm-solidity exists in $(FHEVM_SOLIDITY_PATH)"
	@if [ ! -d $(WORKDIR)/fhevm ]; then \
        echo 'fhevm-solidity is not available in $(WORKDIR)'; \
        echo "FHEVM_SOLIDITY_PATH is set to a custom value"; \
    else \
        echo 'fhevm-solidity is already available in $(WORKDIR)'; \
    fi
else
	@echo "fhevm-solidity does not exist"
	echo "We clone it for you!"
	echo "If you want your own version please update FHEVM_SOLIDITY_PATH pointing to your fhevm-solidity folder!"
	$(MAKE) clone-fhevm-solidity
endif


check-all-test-repo: check-fhevm-solidity

prepare-local-build:
	@bash scripts/replace_go_mod.sh go.mod ".."
	@echo 'go.mod has been updated'
	@git diff go.mod

update-go-mod:
	@cp go.mod $(UPDATE_GO_MOD)
	@bash scripts/replace_go_mod.sh $(UPDATE_GO_MOD) "./work_dir"
	@mv go.mod go.mod.old
	@mv $(UPDATE_GO_MOD) go.mod


build-local:   go.sum $(BUILDDIR)/
	$(info build-local for docker build)
	go build $(BUILD_FLAGS) -o build $(BUILD_ARGS) ./...
	@echo 'ethermintd binary is ready in build folder'

build-local-docker:
ifeq ($(GITHUB_ACTIONS),true)
	$(info Running in a GitHub Actions workflow)
else
	$(info Not running in a GitHub Actions workflow)
	@$(MAKE) clone-go-ethereum
	@$(MAKE) clone-fhevm-go
endif
	$(MAKE) update-go-mod
	@docker compose  -f docker-compose/docker-compose.local.yml build ethermintnodelocal

change-running-node-owner:
ifeq ($(GITHUB_ACTIONS),true)
	$(info Running e2e-test-local in a GitHub Actions workflow)
	sudo chown -R runner:docker running_node/
else
	$(info Not running e2e-test-local in a GitHub Actions workflow)
	@$(SUDO) chown -R $(USER): running_node/
endif


init-ethermint-node:
ifeq ($(LOCAL_BUILD),true)
	$(info LOCAL_BUILD is set, init locally built docker images)
	@$(MAKE) init-ethermint-node-local
else
	$(info LOCAL_BUILD is not set, init docker images from docker registry)
	@$(MAKE) init-ethermint-node-from-registry
endif

init-ethermint-node-local:
	@docker compose -f docker-compose/docker-compose.local.yml run ethermintnodelocal bash /config/setup.sh
	$(MAKE) change-running-node-owner
	$(MAKE) generate-fhe-keys

init-ethermint-node-from-registry:
	@docker compose -f docker-compose/docker-compose.validator.yml run validator bash /config/setup.sh
	$(MAKE) change-running-node-owner
	$(MAKE) generate-fhe-keys-registry

generate-fhe-keys-registry:
ifeq ($(KEY_GEN),false)
	@echo "KEY_GEN is false, executing corresponding commands..."
	@bash ./scripts/copy_fhe_keys.sh $(KMS_DEV_VERSION) $(PWD)/running_node/node2/.ethermintd/zama/keys/network-fhe-keys $(PWD)/running_node/node2/.ethermintd/zama/keys/kms-fhe-keys
else ifeq ($(KEY_GEN),true)
	@echo "KEY_GEN is true, executing corresponding commands..."
	@bash ./scripts/prepare_volumes_from_kms_core.sh $(KMS_DEV_VERSION) $(PWD)/running_node/node2/.ethermintd/zama/keys/network-fhe-keys $(PWD)/running_node/node2/.ethermintd/zama/keys/kms-fhe-keys
else
	@echo "KEY_GEN is set to an unrecognized value: $(KEY_GEN)"
endif

generate-fhe-keys:
ifeq ($(KEY_GEN),false)
	@echo "KEY_GEN is false, executing corresponding commands..."
	@bash ./scripts/copy_fhe_keys.sh $(KMS_DEV_VERSION) $(PWD)/running_node/node1/.ethermintd/zama/keys/network-fhe-keys $(PWD)/running_node/node1/.ethermintd/zama/keys/kms-fhe-keys
else ifeq ($(KEY_GEN),true)
	@echo "KEY_GEN is true, executing corresponding commands..."
	@bash ./scripts/prepare_volumes_from_kms_core.sh $(KMS_DEV_VERSION) $(PWD)/running_node/node1/.ethermintd/zama/keys/network-fhe-keys $(PWD)/running_node/node1/.ethermintd/zama/keys/kms-fhe-keys
else
	@echo "KEY_GEN is set to an unrecognized value: $(KEY_GEN)"
endif

run-ethermint:
ifeq ($(LOCAL_BUILD),true)
	$(info LOCAL_BUILD is set, run locally built docker images)
	@docker compose  -f docker-compose/docker-compose.local.yml -f docker-compose/docker-compose.local.override.yml  up --detach --remove-orphans
else
	$(info LOCAL_BUILD is not set, run docker images from docker registry)
	@docker compose  -f docker-compose/docker-compose.validator.yml -f docker-compose/docker-compose.validator.override.yml  up --detach
endif
	@echo 'sleep a little to let the docker start up'
	sleep 10

stop-ethermint:
ifeq ($(LOCAL_BUILD),true)
	$(info LOCAL_BUILD is set, stop locally built docker images)
	@docker compose  -f docker-compose/docker-compose.local.yml down
else
	$(info LOCAL_BUILD is not set, stop docker images from docker registry)
	@docker compose  -f docker-compose/docker-compose.validator.yml down

endif

TEST_FILE := run_tests.sh
TEST_IF_FROM_REGISTRY := 



install-packages:
	@cd $(FHEVM_SOLIDITY_PATH) && npm ci
	@if [ "$(IS_LINUX)" = "true" ]; then \
	    cd $(FHEVM_SOLIDITY_PATH) && npm i solidity-comments-linux-x64-gnu; \
	fi

run-e2e-test: check-all-test-repo
	$(MAKE) install-packages
# We need to wait so long because as soon as the node start, calling the faucet requires 
# huge amount of gas. In this case the faucet call is failing. So we wait until we reach a
# more stable state which requires less gas. 
	@sleep 10
ifeq ($(LOCAL_BUILD),true)
	@cp ./scripts/run_tests.sh $(FHEVM_SOLIDITY_PATH)/ci/scripts/
	@./fund_test_addresses_docker.sh
	@cd $(FHEVM_SOLIDITY_PATH) && ci/scripts/run_tests.sh
else
	@./fund_test_addresses_docker.sh
	@cd $(FHEVM_SOLIDITY_PATH) && cp .env.example .env
	@cd $(FHEVM_SOLIDITY_PATH) && npm i
	@cd $(FHEVM_SOLIDITY_PATH) && npm run codegen
	@cd $(FHEVM_SOLIDITY_PATH) && ./setup-local-fhevm.sh
	@cd $(FHEVM_SOLIDITY_PATH) && npx hardhat test
endif

e2e-test:
	@$(MAKE) check-all-test-repo
ifeq ($(LOCAL_BUILD),true)
	$(info LOCAL_BUILD is set, run e2e test with locally built docker images)
	$(MAKE) init-ethermint-node-local
else
	$(info LOCAL_BUILD is not set, run e2e test with docker registry for docker images)
	@$(MAKE) init-ethermint-node-from-registry
endif
	$(MAKE) run-ethermint
	$(MAKE) run-e2e-test
	$(MAKE) stop-ethermint


clean-node-storage:
	@echo 'clean node storage'
	sudo rm -rf running_node

clean: clean-node-storage
	$(MAKE) stop-ethermint
	rm -rf $(BUILDDIR)/
	rm -rf $(WORKDIR)/ 
	rm -f $(UPDATE_GO_MOD)


print-info:
	@echo 'GO_ETHEREUM_TAG: $(GO_ETHEREUM_VERSION) ---extracted from go.mod'
	@echo 'FHEVM_GO_TAG: $(FHEVM_GO_VERSION) ---extracted from go.mod'
	@echo 'KMS_DEV_VERSION: $(KMS_DEV_VERSION) for KEY_GEN---extracted from Makefile'
	@echo 'FHEVM_SOLIDITY_VERSION: $(FHEVM_SOLIDITY_VERSION) ---extracted from Makefile'
	@bash scripts/get_repository_info.sh zbc-ethermint ${CURDIR}
	@bash scripts/get_repository_info.sh fhevm $(FHEVM_SOLIDITY_PATH)
	@echo 'AGGREGATOR_TAG: $(shell ./scripts/get_docker_compose_tag.sh docker-compose/docker-compose.local.yml aggregator) --extract from docker-compose.local.yml'
	@echo 'ORACLE_TAG: $(shell ./scripts/get_docker_compose_tag.sh docker-compose/docker-compose.local.yml oracle) --extract from docker-compose.local.yml'
	@echo 'KMS_DEV_VERSION: $(shell ./scripts/get_docker_compose_tag.sh docker-compose/docker-compose.local.yml kms-dev) --extract from docker-compose.local.yml'
	

###############################################################################
###                                Observability                            ###
###############################################################################

observability-start:
	cd observability && docker compose up -d

observability-stop:
	cd observability && docker compose down
