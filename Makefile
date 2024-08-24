.PHONY: yogo tree clean vendor

clean:
	rm -rf ./bin/*

vendor:
	go mod vendor

yogo: clean vendor
	export GOPROXY=https://goproxy.io,direct
	go build -o ./bin/yogo .

tree:
	tree -I node_modules -I vendor -L 2
