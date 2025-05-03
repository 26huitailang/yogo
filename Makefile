.PHONY: yogo tree clean vendor

clean:
	rm -rf ./bin/*

vendor:
	go mod vendor

yogo: clean vendor
	export GOPROXY=https://goproxy.io,direct
	go build -o ./bin/yogo .

install-dev: yogo
	cp ./bin/yogo ~/bin/yogo-dev


install-prod:
	go install github.com/26huitailang/yogo/cmd/yogo@latest

tree:
	tree -I node_modules -I vendor -L 2
