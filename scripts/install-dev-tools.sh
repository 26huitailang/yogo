if ! command -v nvm > /dev/null 2>&1; then
	nvm install --lts
fi
yarn config set registry https://registry.npm.taobao.org/
yarn

go mod vendor
