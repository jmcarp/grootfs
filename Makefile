.PHONY: all \
	test concourse-groot-test concourse-root-test concourse-test \
	go-vet concourse-go-vet go-generate \
	image push-image \
	update-deps

all:
	GOOS=linux go build -o grootfs .
	GOOS=linux go build -o drax ./store/volume_driver/drax

###### Help ###################################################################

help:
	@echo '    all ................................. builds the grootfs cli'
	@echo '    test ................................ runs tests locally'
	@echo '    concourse-groot-test ................ runs groot tests in concourse-lite'
	@echo '    concourse-root-test ................. runs root tests in concourse-lite'
	@echo '    concourse-test ...................... runs tests in concourse-lite'
	@echo '    go-vet .............................. runs go vet in grootfs source code'
	@echo '    concourse-go-vet .................... runs go vet in concourse-lite'
	@echo '    go-generate ......................... runs go generate in grootfs source code'
	@echo '    image ............................... builds a docker image'
	@echo '    push-image .......................... pushes image to docker-hub'
	@echo '    update-deps ......................... update the depedencies'

###### Testing ################################################################

test:
	ginkgo -r -p -race -skipPackage integration .

concourse-test:
	./hack/run-tests -r -g "-p"

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

###### Depedency management ###################################################

update-deps:
	./hack/update-deps
