.PHONY: all \
	concourse-test dracorex-test \
	go-vet concourse-go-vet go-generate \
	image push-image \
	update-deps unit

all:
	GOOS=linux go build -o grootfs .
	GOOS=linux go build -o drax ./store/filesystems/btrfs/drax
	GOOS=linux go build -o tardis ./store/filesystems/overlayxfs/tardis

windows:
	GOOS=windows go build -o grootfs.exe .

###### Help ###################################################################

help:
	@echo '    all ................................. builds the grootfs cli'
	@echo '    windows ............................. builds grootfs for windows'
	@echo '    deps ................................ installs dependencies'
	@echo '    update-deps ......................... updates dependencies'
	@echo '    unit ................................ run unit tests'
	@echo '    concourse-test ...................... runs tests in concourse-lite'
	@echo '    dracorex-test ....................... runs tests on remote CI'
	@echo '    compile-tests ....................... checks that tests can be compiled'
	@echo '    go-vet .............................. runs go vet in grootfs source code'
	@echo '    concourse-go-vet .................... runs go vet in concourse-lite'
	@echo '    go-generate ......................... runs go generate in grootfs source code'
	@echo '    image ............................... builds a docker image'
	@echo '    push-image .......................... pushes image to docker-hub'

###### Dependencies ###########################################################

deps:
	git submodule update --init --recursive

###### Testing ################################################################

compile-tests:
	ginkgo build -r .; find . -name '*.test' | xargs rm

unit:
	./hack/run-tests -r -g "--skipPackage=integration"

concourse-test: go-vet
	./hack/run-tests -r -g "-p"

dracorex-test:
	./hack/run-tests -d -r -g "-p"

###### Go tools ###############################################################

go-vet:
	GOOS=linux go vet `go list ./... | grep -v vendor`

concourse-go-vet:
	fly -t lite e -x -c ci/tasks/go-vet.yml -i grootfs-git-repo=${PWD}

go-generate:
	GOOS=linux go generate `go list ./... | grep -v vendor`

###### Docker #################################################################

image:
	docker build -t cfgarden/grootfs-ci .

push-image:
	docker push cfgarden/grootfs-ci
