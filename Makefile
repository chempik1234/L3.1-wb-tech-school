docker_up:
	docker compose -f docker/docker-compose.yaml up -d

docker_upd:
	docker compose -f docker/docker-compose.yaml up -d

docker_updb:
	docker compose -f docker/docker-compose.yaml up -d --build

docker_down:
	docker compose -f docker/docker-compose.yaml down

lint_delayed_notifier:
	cd delayed_notifier && \
		golint ./... && \
		golangci-lint run ./...

unittest_delayed_notifier:
	cd delayed_notifier && \
		go test ./tests -v --coverprofile=./tests/cover.out --coverpkg=./pkg/pkgports/adapters/cache/lru && \
		go tool cover --html=./tests/cover.out -o ./tests/cover.html

docker_integration_test:
	cd integration_tests && \
		docker compose up -d rabbitmq delayed_notifier simulator_service nginx && \
		docker compose up --build e2e_test
	cd integration_tests && \
		docker compose down