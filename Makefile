
TAG=v0.0.1

all: build

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Build

build: ## Build agent and manager binary.
	go build -o backend/gomonitor-manager backend/backend.go
	go build -o agent/gomonitor-agent agent/agent.go
	go build -o tests/gomonitor-testprog tests/registry.go tests/busy_loop.go tests/test_prog.go

docker-build: ## Build docker image with the agent and manager.
	docker build -t centos:gopacket -f agent/GopacketDockerfile .
	docker build -t gomonitor-manager:${TAG} -f backend/Dockerfile .
	docker build -t gomonitor-agent:${TAG} -f agent/Dockerfile .
	docker build -t gomonitor-testprog:${TAG} -f tests/Dockerfile .

docker-tag: ## Tag docker image for 1445277435
	docker tag gomonitor-manager:${TAG} 1445277435/gomonitor-manager:${TAG}
	docker tag gomonitor-agent:${TAG} 1445277435/gomonitor-agent:${TAG}
	docker tag gomonitor-testprog:${TAG} 1445277435/gomonitor-testprog:${TAG}

docker-push: ## Push docker image with the agent and manager.
	docker push 1445277435/gomonitor-agent:${TAG}
	docker push 1445277435/gomonitor-manager:${TAG}
	docker push 1445277435/gomonitor-testprog:${TAG}

cert: ## generate cert for webhook
	cd backend/config/cert/ && chmod +x generateCert.sh && ./generateCert.sh