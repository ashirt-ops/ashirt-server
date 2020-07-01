
##### Update Targets

.PHONY: update
update: update-go update-frontend

.PHONY: update-go
update-go:
	go mod download

.PHONY: update-frontend
update-frontend:
	cd frontend && npm install

##### Clean Targets

.PHONY: clean-all
clean-all: clean-backend

.PHONY: clean-backend
clean-backend:
	docker-compose down -v

##### Test Targets

.PHONY: test-all
test-all: test-backend test-frontend

.PHONY: test-backend
test-backend:
	cd backend && go test ./...

.PHONY: test-frontend
test-frontend:
	cd frontend && npm run test


##### Linting/Tidy/Format Targets

.PHONY: tidy-all
tidy-all:  format-backend lint-frontend

.PHONY: format-backend
format-backend:
	gofmt -w backend

.PHONY: lint-frontend
lint-frontend:
	cd frontend && npm run lint

# tidy-go removes unused/outdated go modules. This frequently causes conflicts so it is not included in tidy-all
.PHONY: tidy-go
tidy-go:
	go mod tidy

##### Run Targets

# run-backend starts up a fresh instance of the backend -- typically what you want when working on new features
.PHONY: run-backend
run-backend: clean-backend rerun-backend

# rerun-Backend restarts a stopped instance of the backend -- useful if you want to preserve database data between runs
.PHONY: rerun-backend
rerun-backend:
	docker-compose up --build

##### Helpers

# view-backend enters the docker container running the backend
.PHONY: view-backend
view-backend:
	docker exec -it ashirt_backend_1 /bin/sh

# view-backend enters the docker container running the database
.PHONY: view-db
view-db:
	docker exec -it ashirt_db_1 /bin/bash

# new-migration generates a new "migration" (database alteration) when a schema/data change is necessary.
.PHONY: new-migration
new-migration:
	bin/create-migration.sh

# prep is shorthand for formatting and testing. Useful when prepping for a new Pull Request.
.PHONY: prep
prep: tidy-all test-all

# generate-dto-types
.PHONY: generate-dto-types
generate-dto-types: 
	go run backend/dtos/gentypes/generate_typescript_types.go > frontend/src/services/data_sources/dtos.ts
