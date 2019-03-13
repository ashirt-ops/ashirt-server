include .env

.env:
	cp .env_template .env

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
clean-all: clean-termrec clean-backend

.PHONY: clean-termrec
clean-termrec:
	rm -rf dist/termrec/*

.PHONY: clean-backend
clean-backend:
	docker-compose down -v

##### Test Targets

.PHONY: test-all
test-all: test-termrec test-backend test-frontend

.PHONY: test-termrec
test-termrec:
	cd cmd/termrec && go test ./...
	cd termrec && go test ./...

.PHONY: test-backend
test-backend:
	cd backend && go test ./...

.PHONY: test-frontend
test-frontend:
	cd frontend && npm run test


##### Linting/Tidy/Format Targets

.PHONY: tidy-all
tidy-all: format-termrec format-backend lint-frontend

.PHONY: format-termrec
format-termrec:
	gofmt -w cmd/termrec
	gofmt -w termrec

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

##### Build Targets

.PHONY: build-termrec-all
build-termrec-all: build-termrec-linux build-termrec-osx

.PHONY: build-termrec-linux
build-termrec-linux: update-go termrec-copy-static
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/termrec/linux/termrec cmd/termrec/*.go

.PHONY: build-termrec-osx
build-termrec-osx: update-go termrec-copy-static
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o dist/termrec/osx/termrec cmd/termrec/*.go

##### Run Targets

.PHONY: run-termrec
run-termrec: termrec-copy-static
	$(eval export $(shell sed -ne 's/ *#.*$$//; /./ s/=.*$$// p' .env))
	go run cmd/termrec/*.go

# run-backend starts up a fresh instance of the backend -- typically what you want when working on new features
.PHONY: run-backend
run-backend: clean-backend rerun-backend

# rerun-Backend restarts a stopped instance of the backend -- useful if you want to preserve database data between runs
.PHONY: rerun-backend
rerun-backend:
	docker-compose up --build

##### Helpers

.PHONY: termrec-copy-static
termrec-copy-static:
	statik -src=termrec/static -dest=cmd/termrec/config -p=static -f

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
