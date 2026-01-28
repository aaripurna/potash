.PHONY: all

build:
	@mkdir -p _build
	@docker build -t gft:builder .
	@docker container create --name gft-builder gft:builder
	@docker container cp gft-builder:/app ./_build
	@docker container rm gft-builder
	@docker rmi gft:builder
