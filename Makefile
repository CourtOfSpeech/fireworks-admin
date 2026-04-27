.PHONY: run build clean test dev stop restart help ent ent-clean wire wire-clean atlas-diff atlas-apply atlas-status atlas-new atlas-baseline atlas-inspect atlas-validate atlas-hash crud crud-list

APP_NAME := fireworks-admin
MAIN_PATH := cmd/server/main.go
BUILD_DIR := bin
BINARY := $(BUILD_DIR)/$(APP_NAME)
PORT := 1323
ENT_DIR := internal/ent
APP_DIR := internal/app

stop:
	@echo "🛑 Stopping server on port $(PORT)..."
	@lsof -ti :$(PORT) | xargs kill -9 2>/dev/null || echo "No process found on port $(PORT)"

run: stop
	@echo "🚀 Starting server..."
	go run $(MAIN_PATH)

restart: stop run

build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BINARY) $(MAIN_PATH)
	@echo "✅ Build complete: $(BINARY)"

clean:
	@rm -rf $(BUILD_DIR)
	@echo "✅ Clean complete"

test:
	go test -v ./...

dev:
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "❌ air not found. Install with: go install github.com/air-verse/air@latest"; \
		echo "Using go run instead..."; \
		go run $(MAIN_PATH); \
	fi

fmt:
	go fmt ./...

lint:
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "❌ golangci-lint not found. Install from: https://golangci-lint.run/usage/install/"; \
	fi

deps:
	go mod download
	go mod tidy

ent:
	@echo "🔧 Generating Ent code..."
	@cd $(ENT_DIR) && go generate ./...
	@echo "✅ Ent code generated successfully"

ent-clean:
	@echo "🧹 Cleaning Ent generated code..."
	@find $(ENT_DIR) -name "*.go" ! -path "$(ENT_DIR)/schema/*" ! -path "$(ENT_DIR)/hook/*" ! -path "$(ENT_DIR)/intercept/*" ! -name "generate.go" -delete 2>/dev/null || true
	@rm -rf $(ENT_DIR)/enttest $(ENT_DIR)/migrate $(ENT_DIR)/predicate $(ENT_DIR)/runtime $(ENT_DIR)/teltent 2>/dev/null || true
	@echo "✅ Ent generated code cleaned"

wire:
	@echo "🔧 Generating Wire code..."
	@cd $(APP_DIR) && wire
	@echo "✅ Wire code generated successfully"

wire-clean:
	@echo "🧹 Cleaning Wire generated code..."
	@rm -f $(APP_DIR)/wire_gen.go 2>/dev/null || true
	@echo "✅ Wire generated code cleaned"

crud:
	@echo "🔧 CRUD Code Generator"
	@read -p "Enter schema name (e.g., User): " name; \
	if [ -z "$$name" ]; then \
		echo "❌ Schema name is required"; \
		exit 1; \
	fi; \
	read -p "Enter output path (e.g., internal/features/user): " output; \
	if [ -z "$$output" ]; then \
		echo "❌ Output path is required"; \
		exit 1; \
	fi; \
	read -p "Add schema name prefix to filenames? (y/N): " prefix; \
	read -p "Force overwrite? (y/N): " force; \
	prefix_flag=""; \
	force_flag=""; \
	if [ "$$prefix" = "y" ] || [ "$$prefix" = "Y" ]; then \
		prefix_flag="-prefix"; \
	fi; \
	if [ "$$force" = "y" ] || [ "$$force" = "Y" ]; then \
		force_flag="-force"; \
	fi; \
	go run ./cmd/generator -name $$name -output $$output $$prefix_flag $$force_flag

crud-list:
	@echo "📋 Available schemas:"
	@go run ./cmd/generator -list

# Atlas migration commands
atlas-diff:
	@echo "📝 Generating migration diff..."
	@atlas migrate diff --env local
	@echo "✅ Migration diff generated"

atlas-apply:
	@echo "🚀 Applying migrations..."
	@atlas migrate apply --env local
	@echo "✅ Migrations applied"

atlas-status:
	@echo "📊 Migration status..."
	@atlas migrate status --env local

atlas-new:
	@read -p "Enter migration name: " name; \
	atlas migrate new $$name --env local
	@echo "✅ New migration created"

atlas-baseline:
	@echo "🔄 Regenerating baseline..."
	@echo "⚠️  This will archive all existing migrations and create a new baseline"
	@read -p "Enter baseline version (e.g., v1.0.0): " version; \
	echo "Archiving old migrations to migrations/archive/$$version/..."; \
	mkdir -p migrations/archive/$$version 2>/dev/null || true; \
	mv migrations/*.sql migrations/archive/$$version/ 2>/dev/null || true; \
	mv migrations/atlas.sum migrations/archive/$$version/ 2>/dev/null || true; \
	echo "Generating new baseline..."; \
	atlas migrate diff baseline --env local; \
	echo "Marking baseline as applied..."; \
	version_file=$$(ls -t migrations/*.sql 2>/dev/null | head -1 | xargs basename | sed 's/.sql$$//'); \
	if [ -n "$$version_file" ]; then \
		echo "Setting version: $$version_file"; \
		atlas migrate set $$version_file --env local; \
	else \
		echo "⚠️  No migration file found to mark as applied"; \
	fi; \
	echo "✅ Baseline regenerated and old migrations archived to migrations/archive/$$version/"

atlas-inspect:
	@echo "🔍 Inspecting database schema..."
	@atlas schema inspect --env local

atlas-validate:
	@echo "✅ Validating migrations..."
	@atlas migrate validate --env local

atlas-hash:
	@echo "🔢 Recalculating checksums..."
	@atlas migrate hash --env local

help:
	@echo "Available commands:"
	@echo "  make run        - Stop existing server and run the application"
	@echo "  make stop       - Stop the server"
	@echo "  make restart    - Restart the server"
	@echo "  make build      - Build the binary"
	@echo "  make clean      - Remove build artifacts"
	@echo "  make test       - Run tests"
	@echo "  make dev        - Run with hot reload (requires air)"
	@echo "  make fmt        - Format code"
	@echo "  make lint       - Run linter"
	@echo "  make deps       - Download dependencies"
	@echo "  make ent        - Generate Ent ORM code"
	@echo "  make ent-clean  - Clean Ent generated code"
	@echo "  make wire       - Generate Wire dependency injection code"
	@echo "  make wire-clean - Clean Wire generated code"
	@echo "  make crud       - Generate CRUD code from Ent schema (interactive)"
	@echo "  make crud-list  - List available schemas for CRUD generation"
	@echo ""
	@echo "Atlas migration commands:"
	@echo "  make atlas-diff     - Generate migration diff from schema"
	@echo "  make atlas-apply    - Apply pending migrations"
	@echo "  make atlas-status   - Show migration status"
	@echo "  make atlas-new      - Create new empty migration"
	@echo "  make atlas-baseline - Regenerate baseline (archive old migrations)"
	@echo "  make atlas-inspect  - Inspect database schema"
	@echo "  make atlas-validate - Validate migration files"
	@echo "  make atlas-hash     - Recalculate migration checksums"
