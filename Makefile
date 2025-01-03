NODE_BIN := ./node_modules/.bin

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^//'

.PHONY: confirm
confirm:
	@echo 'Are you sure? [y/N]' && read ans && [ $${ans:-N} = y ]

## build: build the application
.PHONY:build
build: gen/css gen/templ
	@go build -tags dev -o bin/htmx-chat ./cmd/server

## build/prod: build the application for production
.PHONY:build/prod
build/prod: gen/css gen/templ
	@go build -tags prod -o bin/htmx-chat ./cmd/server

## run: run the application
.PHONY:run
run: build
	@./bin/htmx-chat

## install: install dependencies
.PHONY: install
install:
	@make install/templ
	@make gen/templ
	@make install/go
	@make install/htmx
	@make install/css

## install/go: install go dependencies
.PHONY: install/go
install/go:
	@go get ./...
	@go mod tidy
	@go mod download

## install/templ: install templ
.PHONY: install/templ
install/templ:
	@go install github.com/a-h/templ/cmd/templ@latest

## install/htmx: install htmx
.PHONY: install/htmx
install/htmx:
	@mkdir -p cmd/server/deps
	@wget -q -O cmd/server/deps/htmx-1.9.12.min.js.gz https://unpkg.com/htmx.org@1.9.12/dist/htmx.min.js.gz
	@gunzip -f cmd/server/deps/htmx-1.9.12.min.js.gz
	@wget -q -O cmd/server/deps/sse-1.9.12.js https://unpkg.com/htmx.org@1.9.12/dist/ext/sse.js
	@go install github.com/tdewolff/minify/v2/cmd/minify@latest
	@minify -o cmd/server/deps/sse-1.9.12.min.js cmd/server/deps/sse-1.9.12.js
	@rm cmd/server/deps/sse-1.9.12.js

## install/css: install css
.PHONY: install/css
install/css:
	@npm install -D tailwindcss

## gen/css: generate css
.PHONY: gen/css
gen/css:
	@$(NODE_BIN)/tailwindcss build -i internal/views/css/app.css -o cmd/server/public/styles.css --minify

## gen/templ: generate templ
.PHONY: gen/templ
gen/templ:
	@templ generate

## watch/css: watch css
.PHONY: watch/css
watch/css:
	@$(NODE_BIN)/tailwindcss -i internal/views/css/app.css -o cmd/server/public/styles.css --minify --watch

## watch/templ: watch templ
.PHONY: watch/templ
watch/templ:
	@templ generate --watch --proxy=http://127.0.0.1:8000

## ci/scaffold: scaffold the project
.PHONY: ci/scaffold
ci/scaffold:
	@mkdir -p cmd/server/deps
	@echo "hello world" > cmd/server/deps/hello.txt
	@mkdir -p cmd/server/public
	@echo "hello world" > cmd/server/public/hello.txt
