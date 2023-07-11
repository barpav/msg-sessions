up:
	sudo docker-compose up -d --wait
down:
	sudo docker-compose down

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