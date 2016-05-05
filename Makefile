# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: gvec gvec-cross evm all test travis-test-with-coverage xgo clean
.PHONY: gvec-linux gvec-linux-arm gvec-linux-386 gvec-linux-amd64
.PHONY: gvec-darwin gvec-darwin-386 gvec-darwin-amd64
.PHONY: gvec-windows gvec-windows-386 gvec-windows-amd64
.PHONY: gvec-android gvec-android-16 gvec-android-21

GOBIN = build/bin

MODE ?= default
GO ?= latest

gvec:
	build/env.sh go install -v $(shell build/flags.sh) ./cmd/gvec
	@echo "Done building."
	@echo "Run \"$(GOBIN)/gvec\" to launch gvec."

gvec-cross: gvec-linux gvec-darwin gvec-windows gvec-android
	@echo "Full cross compilation done:"
	@ls -l $(GOBIN)/gvec-*

gvec-linux: xgo gvec-linux-arm gvec-linux-386 gvec-linux-amd64
	@echo "Linux cross compilation done:"
	@ls -l $(GOBIN)/gvec-linux-*

gvec-linux-386: xgo
	build/env.sh $(GOBIN)/xgo --go=$(GO) --buildmode=$(MODE) --dest=$(GOBIN) --targets=linux/386 -v $(shell build/flags.sh) ./cmd/gvec
	@echo "Linux 386 cross compilation done:"
	@ls -l $(GOBIN)/gvec-linux-* | grep 386

gvec-linux-amd64: xgo
	build/env.sh $(GOBIN)/xgo --go=$(GO) --buildmode=$(MODE) --dest=$(GOBIN) --targets=linux/amd64 -v $(shell build/flags.sh) ./cmd/gvec
	@echo "Linux amd64 cross compilation done:"
	@ls -l $(GOBIN)/gvec-linux-* | grep amd64

gvec-linux-arm: gvec-linux-arm-5 gvec-linux-arm-6 gvec-linux-arm-7 gvec-linux-arm64
	@echo "Linux ARM cross compilation done:"
	@ls -l $(GOBIN)/gvec-linux-* | grep arm

gvec-linux-arm-5: xgo
	build/env.sh $(GOBIN)/xgo --go=$(GO) --buildmode=$(MODE) --dest=$(GOBIN) --targets=linux/arm-5 -v $(shell build/flags.sh) ./cmd/gvec
	@echo "Linux ARMv5 cross compilation done:"
	@ls -l $(GOBIN)/gvec-linux-* | grep arm-5

gvec-linux-arm-6: xgo
	build/env.sh $(GOBIN)/xgo --go=$(GO) --buildmode=$(MODE) --dest=$(GOBIN) --targets=linux/arm-6 -v $(shell build/flags.sh) ./cmd/gvec
	@echo "Linux ARMv6 cross compilation done:"
	@ls -l $(GOBIN)/gvec-linux-* | grep arm-6

gvec-linux-arm-7: xgo
	build/env.sh $(GOBIN)/xgo --go=$(GO) --buildmode=$(MODE) --dest=$(GOBIN) --targets=linux/arm-7 -v $(shell build/flags.sh) ./cmd/gvec
	@echo "Linux ARMv7 cross compilation done:"
	@ls -l $(GOBIN)/gvec-linux-* | grep arm-7

gvec-linux-arm64: xgo
	build/env.sh $(GOBIN)/xgo --go=$(GO) --buildmode=$(MODE) --dest=$(GOBIN) --targets=linux/arm64 -v $(shell build/flags.sh) ./cmd/gvec
	@echo "Linux ARM64 cross compilation done:"
	@ls -l $(GOBIN)/gvec-linux-* | grep arm64

gvec-darwin: gvec-darwin-386 gvec-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -l $(GOBIN)/gvec-darwin-*

gvec-darwin-386: xgo
	build/env.sh $(GOBIN)/xgo --go=$(GO) --buildmode=$(MODE) --dest=$(GOBIN) --targets=darwin/386 -v $(shell build/flags.sh) ./cmd/gvec
	@echo "Darwin 386 cross compilation done:"
	@ls -l $(GOBIN)/gvec-darwin-* | grep 386

gvec-darwin-amd64: xgo
	build/env.sh $(GOBIN)/xgo --go=$(GO) --buildmode=$(MODE) --dest=$(GOBIN) --targets=darwin/amd64 -v $(shell build/flags.sh) ./cmd/gvec
	@echo "Darwin amd64 cross compilation done:"
	@ls -l $(GOBIN)/gvec-darwin-* | grep amd64

gvec-windows: xgo gvec-windows-386 gvec-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -l $(GOBIN)/gvec-windows-*

gvec-windows-386: xgo
	build/env.sh $(GOBIN)/xgo --go=$(GO) --buildmode=$(MODE) --dest=$(GOBIN) --targets=windows/386 -v $(shell build/flags.sh) ./cmd/gvec
	@echo "Windows 386 cross compilation done:"
	@ls -l $(GOBIN)/gvec-windows-* | grep 386

gvec-windows-amd64: xgo
	build/env.sh $(GOBIN)/xgo --go=$(GO) --buildmode=$(MODE) --dest=$(GOBIN) --targets=windows/amd64 -v $(shell build/flags.sh) ./cmd/gvec
	@echo "Windows amd64 cross compilation done:"
	@ls -l $(GOBIN)/gvec-windows-* | grep amd64

gvec-android: xgo
	build/env.sh $(GOBIN)/xgo --go=$(GO) --buildmode=$(MODE) --dest=$(GOBIN) --targets=android/* -v $(shell build/flags.sh) ./cmd/gvec
	@echo "Android cross compilation done:"
	@ls -l $(GOBIN)/gvec-android-*

gvec-ios: gvec-ios-arm-7 gvec-ios-arm64
	@echo "iOS cross compilation done:"
	@ls -l $(GOBIN)/gvec-ios-*

gvec-ios-arm-7: xgo
	build/env.sh $(GOBIN)/xgo --go=$(GO) --buildmode=$(MODE) --dest=$(GOBIN) --targets=ios/arm-7 -v $(shell build/flags.sh) ./cmd/gvec
	@echo "iOS ARMv7 cross compilation done:"
	@ls -l $(GOBIN)/gvec-ios-* | grep arm-7

gvec-ios-arm64: xgo
	build/env.sh $(GOBIN)/xgo --go=$(GO) --buildmode=$(MODE) --dest=$(GOBIN) --targets=ios-7.0/arm64 -v $(shell build/flags.sh) ./cmd/gvec
	@echo "iOS ARM64 cross compilation done:"
	@ls -l $(GOBIN)/gvec-ios-* | grep arm64

evm:
	build/env.sh $(GOROOT)/bin/go install -v $(shell build/flags.sh) ./cmd/evm
	@echo "Done building."
	@echo "Run \"$(GOBIN)/evm to start the evm."

all:
	build/env.sh go install -v $(shell build/flags.sh) ./...

test: all
	build/env.sh go test ./...

travis-test-with-coverage: all
	build/env.sh build/test-global-coverage.sh

xgo:
	build/env.sh go get github.com/karalabe/xgo

clean:
	rm -fr build/_workspace/pkg/ Godeps/_workspace/pkg $(GOBIN)/*
