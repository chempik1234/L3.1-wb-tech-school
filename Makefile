docker_up:
	docker compose -f docker/docker-compose.yaml up -d

docker_upd:
	docker compose -f docker/docker-compose.yaml up -d

docker_updb: build docker_upd

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
		docker compose up -d rabbitmq delayed_notifier nginx && \
		docker compose up --build e2e_test
	cd integration_tests && \
		docker compose down

build:
	docker build -t delayed_notifier -f docker/service.Dockerfile ./delayed_notifier

kubernetes_up:
	echo "kubectl apply -f ?"
	# TODO: kubernetes up
	# TODO: kubernetes down
	# TODO: kubernetes apply

init_env:
	type ".\config\example.env" > ".\config\.env"
