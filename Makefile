.PHONY: start
start:
	docker-compose build
	docker-compose up

.PHONY: stop
stop:
	docker-compose stop -t 0
	docker-compose rm -f
