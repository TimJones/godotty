PROJECT=godotty
CODE_DIRS=cmd/ internal/
IMAGE=quay.io/timjones/$(PROJECT)
DOCKER:=$(shell command -v docker 2>/dev/null)
DOCKER_COMPOSE:=$(shell command -v docker-compose 2>/dev/null)
UID:=$(shell id -u)
GID:=$(shell id -g)
.DEFAULT_GOAL:=bin

ifdef NO_DOCKER
define docker_wrapper
	$(1)
endef
else
define docker_wrapper
	@echo $(1)
	@$(DOCKER_COMPOSE) --project-name $(PROJECT) run \
		--rm \
		--user "$(UID):$(GID)" \
		$(PROJECT) \
		$(1)
endef
endif

.PHONY: shell
shell:
	$(call docker_wrapper,bash)

.PHONY: bin
bin:
	$(call docker_wrapper,go build -a -ldflags '-extldflags "-static"' -o bin/$(PROJECT) ./cmd/$(PROJECT))

.PHONY: check-fmt
check-fmt:
	$(call docker_wrapper,sh -c 'gofmt -s -l -e -d $(CODE_DIRS) && exit `gofmt -s -l $(CODE_DIRS) | wc -l`')

.PHONY: fmt
fmt:
	$(call docker_wrapper,gofmt -s -l -w $(CODE_DIRS))

.PHONY: test
test:
	$(call docker_wrapper,go test ./...)
