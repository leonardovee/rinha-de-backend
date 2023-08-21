run:
	docker container prune && docker volume prune && docker-compose up --build
