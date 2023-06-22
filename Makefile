up:
	sudo docker-compose up -d --wait
down:
	sudo docker-compose down

run-r:
	sudo docker rm -f msg-storage-sessions-v1
	sudo docker run --name msg-storage-sessions-v1 -p 6379:6379 -d redis:alpine3.18
	sudo docker ps
exec-r:
	sudo docker exec -it msg-storage-sessions-v1 redis-cli
stop-r:
	sudo docker stop msg-storage-sessions-v1
	sudo docker rm -f msg-storage-sessions-v1
	sudo docker ps