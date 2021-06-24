
TAG=v0.0.1

all: build

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Build

build: ## Build agent and manager binary.
	go build -o bin/kube-gomonitor-backend backend/backend.go
	go build -o bin/kube-gomonitor-agent agent/agent.go
	go build -o bin/kube-gomonitor-testprog tests/registry.go tests/busy_loop.go tests/test_prog.go

docker-build: ## Build docker image with the agent and manager.
	cd agent && docker build -t 1445277435/centos:gopacket -f docker/Dockerfile.Gopacket ..
	cd agent && docker build -t 1445277435/kube-gomonitor-agent:${TAG} -f docker/Dockerfile ..
	cd backend && docker build -t 1445277435/kube-gomonitor-backend:${TAG} -f docker/Dockerfile ..
	cd tests && docker build -t 1445277435/kube-gomonitor-testprog:${TAG} -f docker/Dockerfile ..

docker-push: ## Push docker image with the agent and manager.
	docker push 1445277435/kube-gomonitor-agent:${TAG}
	docker push 1445277435/kube-gomonitor-backend:${TAG}
	docker push 1445277435/kube-gomonitor-testprog:${TAG}

cert: ## generate cert for webhook
	cd backend/deployment/cert/ && chmod +x generateCert.sh && ./generateCert.sh