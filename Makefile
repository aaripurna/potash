.PHONY: all build dev

# Prefer docker, fall back to podman. Override with: make build CONTAINER=podman
CONTAINER ?= $(shell for c in docker podman; do command -v $$c >/dev/null 2>&1 && { echo $$c; break; }; done)

build:
	@test -n "$(CONTAINER)" || { echo "make: no container runtime found - install docker or podman"; exit 1; }
	@echo "==> using $(CONTAINER)"
	@mkdir -p _build
	@$(CONTAINER) build -t gft:builder .
	@$(CONTAINER) container create --name gft-builder gft:builder
	@$(CONTAINER) container cp gft-builder:/app ./_build
	@$(CONTAINER) container rm gft-builder
	@$(CONTAINER) rmi gft:builder

dev:
	@foreman s -f Procfile
