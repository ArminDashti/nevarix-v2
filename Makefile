.PHONY: proto build-hub build-agent build-manager test compose-up compose-down certs

ROOT := $(CURDIR)
HUB_DIR := apps/hub
AGENT_DIR := apps/agent
MANAGER_DIR := apps/manager
PROTO_DIR := packages/proto

proto:
	cd $(PROTO_DIR) && buf generate

build-hub:
	cd $(HUB_DIR) && go build -o ../../bin/nvx-hub ./cmd/nvx-hub

build-agent:
	cd $(AGENT_DIR) && go build -o ../../bin/nvx-agent ./cmd/nvx-agent

build-manager:
	cd $(MANAGER_DIR) && npm run build

test:
	cd $(HUB_DIR) && go test ./...
	cd $(AGENT_DIR) && go test ./...
	cd $(MANAGER_DIR) && npm test --if-present

compose-up:
	docker compose -f deploy/docker-compose.dev.yml up -d --build

compose-down:
	docker compose -f deploy/docker-compose.dev.yml down

certs:
	bash scripts/dev/generate-certs.sh

build: build-hub build-agent

all: proto build
