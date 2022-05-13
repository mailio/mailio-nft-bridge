.PHONY: build-staging
build-staging: ## Build the staging docker image.
	docker compose -f docker/staging/docker-compose.yml build

.PHONY: start-staging
start-staging: ## Start the staging docker container.
	docker compose -f docker/staging/docker-compose.yml up -d

.PHONY: stop-staging
stop-staging: ## Stop the staging docker container.
	docker compose -f docker/staging/docker-compose.yml down

.PHONY: make-user-linux-amd64
make-user-linux-amd64: ## Build the make user script command
	OOS=linux GOARCH=amd64 go build -o ./scripts/make-user ./scripts/make_user.go