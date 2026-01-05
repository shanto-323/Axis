run:
	go run ./cmd/main.go

migration_new: 
	@if [ -z "$(name)" ]; then \
		@echo "Error: name parameter is required" \
        @echo "Usage: task migrations:new name=migration_name" \
		exit 1; \
	fi
	@echo "Creating migration file for $(name)"
	@tern new -m ./internal/database/migrations ${name}

migration_up: confirm
	@if [ -z "$(DSN)" ]; then \
		echo "Error: DSN is required"; \
		echo "Usage: make migrations-up DSN=\"postgres://user:pass@host:5432/dbname?sslmode=disable\""; \
		exit 1; \
	fi
	@echo "Running up migrations..."
	tern migrate -m ./internal/database/migrations --conn-string "$(DSN)"

confirm:
	@read -p "Press [y/Y] to continue: " value; \
	value=$$(echo $$value | tr '[:upper:]' '[:lower:]'); \
	if [ "$$value" = "y" ]; then \
		echo "Running target..."; \
	else \
		exit 1; \
	fi
