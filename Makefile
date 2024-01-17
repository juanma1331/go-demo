# For Windows
# If you are using Linux or Mac, please adapt the commands accordingly (e.g. basename, rm, ...)

PACKAGES := $(shell go list ./...)
name := $(shell powershell -Command "(Get-Location).Path.Split('\\')[-1]")

.PHONY: print-vars
print-vars:
	@echo $(PACKAGES)
	@echo $(name)

.PHONY: vet
vet:
	@go vet $(PACKAGES)


.PHONY: test
test:
	@powershell -Command "$$env:CGO_ENABLED='1'; go test -race -cover $(PACKAGES)"


.PHONY: build-dev
build-dev:
	@templ generate
	@go build -o ./tmp/main-dev.exe ./cmd/main.go


.PHONY: start
start: build-dev
	@air -c air.toml


# Download tailwindcss binary first
.PHONY: css
css:
	@./bin/tailwindcss -i ./assets/css/input.css -o ./assets/css/output.css --minify


.PHONY: css-watch
css-watch:
	@./bin/tailwindcss -i ./assets/css/input.css -o ./assets/css/output.css --watch


.PHONY: migrate-db
migrate-db:
	@go run ./cmd/migration/main.go


.PHONY: seed-db
seed-db:
	@go run ./cmd/seed/main.go


.PHONY: delete-db
delete-db:
	@del .\data\data.sqlite


.PHONY: reset-db
reset-db: delete-db migrate-db seed-db