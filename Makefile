PWD = $(shell pwd)

lint:
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint golangci-lint run -v

up:
	docker-compose -f "deployment/docker-compose.yml" up -d --build

down:
	docker-compose -f "deployment/docker-compose.yml" down