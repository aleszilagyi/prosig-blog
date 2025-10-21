.PHONY: format
format:
	find . -type f -name '*.go' ! -name 'mock_*.go' -exec gofmt -s -w {} +

.PHONY: test
test:
	go test ./... -race -shuffle=on

.PHONY: local-build
local-build: build-dir
	go build -o _build/main ./cmd

.PHONY: build-dir
build-dir:
	mkdir -p _build

.PHONY: mocks
mocks: install-codegen
	@echo "Checking for mockgen..."
	@if [ -x "$$HOME/go/bin/mockgen" ]; then \
		MOCKGEN_PATH="$$HOME/go/bin"; \
	elif [ -n "$$GOPATH" ] && [ -x "$$GOPATH/bin/mockgen" ]; then \
		MOCKGEN_PATH="$$GOPATH/bin"; \
	else \
		echo "Error: mockgen not found in $$HOME/go/bin or $$GOPATH/bin"; \
		echo "Install it with: go install github.com/golang/mock/mockgen@latest"; \
		exit 1; \
	fi; \
	echo "Running go generate with mockgen"; \
	PATH="$$PATH:$$MOCKGEN_PATH" go generate ./...

.PHONY: run-docker-local
run-docker-local:
	docker compose up -d --build

.PHONY: down-docker-local
down-docker-local:
	docker compose down

.PHONY: install-codegen
install-codegen:
	@echo "--- Installing mock gen tools... ---"
	@sh scripts/install_codegen_tools.sh

.PHONY: request-get-posts
request-get-posts:
	@curl -X GET http://localhost:8080/api/posts

.PHONY: request-get-post-1
request-get-post-1:
	@curl -X GET http://localhost:8080/api/posts/1

.PHONY: request-get-post-fail
request-get-post-fail:
	@curl -X GET http://localhost:8080/api/posts/a

.PHONY: request-post-comment-post-1
request-post-comment-post-1:
	@curl -X POST http://localhost:8080/api/posts/1/comments \
			-H "Content-Type: application/json" \
			-d '{"comment_content": "Great post!"}'

.PHONY: request-post-comment-post-validation-fail
request-post-comment-post-validation-fail:
	@curl -X POST http://localhost:8080/api/posts/1/comments \
			-H "Content-Type: application/json" \
			-d '{"comment_content": ""}'

.PHONY: request-post-comment-post-postid-fail
request-post-comment-post-postid-fail:
	@curl -X POST http://localhost:8080/api/posts/a/comments \
			-H "Content-Type: application/json" \
			-d '{"comment_content": "Great Post"}'

.PHONY: request-post-post
request-post-post:
	@curl -X POST http://localhost:8080/api/posts \
			-H "Content-Type: application/json" \
			-d '{"title": "My first post", "post_content": "Hello world!"}'

.PHONY: request-post-post-fail
request-post-post-fail:
	@curl -X POST http://localhost:8080/api/posts \
			-H "Content-Type: application/json" \
			-d '{"title": "", "post_content": "Hello world!"}'
