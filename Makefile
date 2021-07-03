.PHONY: start
start:
	docker-compose build
	docker-compose up

.PHONY: stop
stop:
	docker-compose stop -t 0
	docker-compose rm -f

.PHONY: install
install:
	echo "치할것들 정리"

.PHONY: generate
generate:
	echo "protoc  및 proto 파일 빌드"


#
# 여기서 protoc 만들고, 기타 설치하는 내용들 정리하자.