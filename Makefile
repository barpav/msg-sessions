up:
	sudo docker-compose up -d --wait
down:
	sudo docker-compose down

build:
	sudo docker image rm -f ghcr.io/barpav/msg-sessions:v1
	sudo docker build -t ghcr.io/barpav/msg-sessions:v1 -f docker/service/Dockerfile .
	sudo docker image ls
clear:
	sudo docker image rm -f ghcr.io/barpav/msg-sessions:v1

up-debug:
	sudo docker-compose -f compose-debug.yaml up -d --wait
down-debug:
	sudo docker-compose -f compose-debug.yaml down

exec-redis:
	sudo docker exec -it msg-storage-sessions-v1 redis-cli

push:
	sudo docker push ghcr.io/barpav/msg-sessions:v1

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    sessions_service_go_grpc/sessions_service.proto