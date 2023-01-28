build:
	docker image build -f Dockerfile -t forum-docker .
run:
	docker container run -p 8080:8080 --detach --name forum-container forum-docker
allstop:
	docker stop $(docker ps -a -q)
prune:
	docker system prune -a
stop:
	docker stop forum-container
	docker ps -a
go-run:
	go run ./cmd/main.go