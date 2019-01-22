PROJECT=godotty
IMAGE=quay.io/timjones/$(PROJECT)
DOCKER:=$(shell command -v docker 2>/dev/null)
DOCKER_COMPOSE:=$(shell command -v docker-compose 2>/dev/null)
UID:=$(shell id -u)
GID:=$(shell id -g)
.DEFAULT_GOAL:=bin

.PHONY: docker-check
docker-check:
ifndef DOCKER
	@echo "docker is not available. Please install docker"
endif
ifndef DOCKER_COMPOSE
	@echo "docker-compose is not available. Please install docker-compose"
endif

.PHONY: dev-container
dev-container: docker-check
	-$(DOCKER) tag $(PROJECT)_$(PROJECT):latest $(PROJECT)_$(PROJECT):old
	$(DOCKER_COMPOSE) --project-name $(PROJECT) build
	-$(DOCKER) rmi $(PROJECT)_$(PROJECT):old

.PHONY: shell
shell: dev-container
	$(DOCKER_COMPOSE) --project-name $(PROJECT) run \
		--name $(PROJECT)-shell \
		--rm \
		--user "$(UID):$(GID)" \
		$(PROJECT) \
		bash

.PHONY: clean
clean: docker-check
	$(DOCKER_COMPOSE) --project-name $(PROJECT) down \
	--rmi local \
	--remove-orphans

.PHONY: image
image: docker-check
	$(DOCKER) build --tag $(IMAGE) .

.PHONY: bin
bin: dev-container
	$(DOCKER_COMPOSE) --project-name $(PROJECT) run \
		--name $(PROJECT)-bin \
		--rm \
		--user "$(UID):$(GID)" \
		$(PROJECT) \
		go build -a -ldflags '-extldflags "-static"' -o bin/$(PROJECT) cmd/$(PROJECT)/*.go

.PHONY: fmt
fmt: dev-container
	$(DOCKER_COMPOSE) --project-name $(PROJECT) run \
		--name $(PROJECT)-fmt \
		--rm \
		--user "$(UID):$(GID)" \
		$(PROJECT) \
		gofmt -s -l -w cmd/ pkg/
