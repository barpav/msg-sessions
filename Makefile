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

user:
	curl -v -X POST	-H "Content-Type: application/vnd.newUser.v1+json" \
	-d '{"id": "jane", "name": "Jane Doe", "password": "My1stGoodPassword"}' \
	localhost:8081
session:
	curl -v -X POST -H "Authorization: Basic amFuZTpNeTFzdEdvb2RQYXNzd29yZA==" localhost:8080
get-sessions:
	curl -v -H "Authorization: Basic amFuZTpNeTFzdEdvb2RQYXNzd29yZA==" \
	-H "Accept: application/vnd.userSessions.v1+json" \
	localhost:8080
end-sessions-grpc:
	grpcurl -import-path sessions_service_go_grpc \
	-proto sessions_service.proto \
	-d '{"id": "jane"}' \
	-plaintext localhost:9000 msg.sessions.Sessions/EndAll

exec-redis:
	sudo docker exec -it msg-storage-sessions-v1 redis-cli

push:
	sudo docker push ghcr.io/barpav/msg-sessions:v1

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    sessions_service_go_grpc/sessions_service.proto