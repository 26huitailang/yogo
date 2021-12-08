tag=latest
build:
	docker build -f docker/Dockerfile -t yogo:$(tag) .