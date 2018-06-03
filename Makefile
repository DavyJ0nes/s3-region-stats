#### VARIABLES ####
USERNAME = davyj0nes
APP_NAME = s3-region-stats

GO_VERSION ?= 1.10.2
GO_PROJECT_PATH ?= github.com/davyj0nes/s3-region-stats
GO_FILES = $(shell go list ./... | grep -v /vendor/)

RELEASE = 0.0.1
COMMIT = $(shell git rev-parse HEAD | cut -c 1-6)
BUILD_TIME = $(shell date -u '+%Y-%m-%d_%I:%M:%S%p')

BUILD_PREFIX = CGO_ENABLED=0 GOOS=linux
BUILD_FLAGS = -a -tags netgo --installsuffix netgo
DOCKER_GO_BUILD = docker run --rm -v "$(GOPATH)":/go -v "$(CURDIR)":/go/src/app -w /go/src/app golang:${GO_VERSION}
GO_BUILD_STATIC = $(BUILD_PREFIX) go build $(BUILD_FLAGS)
GO_BUILD_OSX = GOOS=darwin GOARCh=amd64 go build
GO_BUILD_WIN = GOOS=windows GOARCh=amd64 go build

DOCKER_RUN_CMD = docker run -it --rm -v ${HOME}/.aws:/root/.aws --name ${APP_NAME} ${USERNAME}/${APP_NAME}:${IMAGE_VERSION} "\$$@"

#### COMMANDS ####
.PHONY: compile
compile:
	@mkdir -p releases/${RELEASE}
	$(call blue, "# Compiling Linux App...")
	@${DOCKER_GO_BUILD} sh -c 'go get && ${GO_BUILD_STATIC} -o releases/${RELEASE}/${APP_NAME}_linux'
	$(call blue, "# Compiling OSX App...")
	@${DOCKER_GO_BUILD} sh -c 'go get && ${GO_BUILD_OSX} -o releases/${RELEASE}/${APP_NAME}_osx'
	$(call blue, "# Compiling Windows App...")
	@${DOCKER_GO_BUILD} sh -c 'go get && ${GO_BUILD_WIN} -o releases/${RELEASE}/${APP_NAME}.exe'
	@$(MAKE) clean

.PHONY: binary
binary:
	$(call blue, "# Building Golang Binary...")
	@${DOCKER_GO_BUILD} sh -c 'go get && ${GO_BUILD_STATIC} -o ${APP_NAME}_static'

.PHONY: image
image: binary
	$(call blue, "# Building Docker Image...")
	@docker build --no-cache --label APP_VERSION=${RELEASE} --label BUILT_ON=${BUILD_TIME} --label GIT_HASH=${COMMIT} -t ${USERNAME}/${APP_NAME}:${RELEASE} .
	@docker tag ${USERNAME}/${APP_NAME}:${RELEASE} ${USERNAME}/${APP_NAME}:latest
	@$(MAKE) clean

# .PHONY: publish
# publish: image
#         $(call blue, "# Publishing Docker Image...")
#         @docker push docker.io/${USERNAME}/${APP}:${RELEASE}

.PHONY: run
run:
	$(call blue, "# Running App...")
	@docker run -it --rm -e "AWS_PROFILE=$(AWS_PROFILE)" -v "$(GOPATH):/go" -v "$(CURDIR):/go/src/app" -v "$(HOME)/.aws:/root/.aws/" -w /go/src/app golang:${GO_VERSION} go run main.go

.PHONY: run_image
run_image: 
	$(call blue, "# Running Docker Image Locally...")
	@docker run -it --rm --name ${APP_NAME} -e "AWS_PROFILE=$(AWS_PROFILE)" -v "$(HOME)/.aws:/home/dockmaster/.aws/" ${USERNAME}/${APP_NAME}:${RELEASE} 

.PHONY: test
test:
	$(call blue, "# Testing Golang Code...")
	@docker run --rm -it -v "$(GOPATH):/go" -v "$(CURDIR)":/go/src/app -w /go/src/app golang:${GO_VERSION} sh -c 'go test -v -race ${GO_FILES}' 

.PHONY: clean
clean: 
	@rm -f ${APP_NAME}_static

#### FUNCTIONS ####
define blue
	@tput setaf 4
	@echo $1
	@tput sgr0
endef
