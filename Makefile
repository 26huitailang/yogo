tag=latest
.PHONY: build tree
build:
	docker build -f docker/Dockerfile -t yogo:$(tag) .

tree:
	tree -I node_modules -I vendor -L 2
