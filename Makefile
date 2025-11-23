.PHONY: reset build up

reset:
	docker-compose down --rmi all -v
	docker volume prune -f
	docker image prune -af

build:
	docker-compose build --no-cache

up:
	docker-compose build --no-cache
	docker-compose up -d
