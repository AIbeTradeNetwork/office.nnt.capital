.PHONY: build build-prof build-run dc mocks gen test run dev lint

build:
	go build -o ./build/server ./cmd/server/main.go

build-partner-processor:
	go build -o ./build/partner-processor ./cmd/partner_applications_processor/main.go

run-partner-processor: build-partner-processor
	./build/partner-processor

build-prof: build
	go tool pprof —text ./build/server

build-run: build
	./build/server

dc:
	docker-compose up  --remove-orphans --build

gen:
	go generate ./...

test:
	go test -v -coverprofile cover.out ./... && go tool cover -html=cover.out

run:
	go run -race ./cmd/main.go

dev:
	air -c api.air.toml

lint:
	golangci-lint run