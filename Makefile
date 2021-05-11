.PHONY: app-create
app-create:
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o ./cmd/app ./cmd
	@docker build -t aig ./cmd

.PHONY: elastic-down
elastic-down:
	@docker-compose -f cmd/docker-compose.yaml down

.PHONY: elastic-up
elastic-up:
	@docker-compose -f cmd/docker-compose.yaml up -d

.PHONY: http-test
http-test:
	@go run ./tools/http/main.go -case=login.json