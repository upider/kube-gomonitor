
TAG=v0.0.1

all: build docker-build docker-push

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Build

build: ## Build agent and manager binary.
	go build -o backend/gomonitor-manager backend/backend.go
	go build -o agent/gomonitor-agent agent/agent.go

docker-build: ## Build docker image with the agent and manager.
	docker build -t centos:gopacket -f agent/GopacketDockerfile .
	docker build -t gomonitor-agent:${TAG} -f agent/Dockerfile .
	docker build -t gomonitor-manager:${TAG} -f backend/Dockerfile .

docker-push: ## Push docker image with the agent and manager.
	docker push gomonitor-agent:${TAG}
	docker push gomonitor-manager:${TAG}