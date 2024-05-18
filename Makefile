docker-compose-up:
	docker-compose -f .docker-compose/docker-compose.yaml up

docker-compose-rebuild:
	docker-compose -f .docker-compose/docker-compose.yaml up -d --no-deps --build

docker-compose-down:
	docker-compose -f .docker-compose/docker-compose.yaml down