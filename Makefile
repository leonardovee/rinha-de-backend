run:
	docker container prune && docker volume prune && docker-compose up --build

stress:
	bash ./stress-test/run-test.sh
