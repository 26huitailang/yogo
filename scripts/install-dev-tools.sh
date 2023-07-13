if ! command -v yarn &> /dev/null; then
	nvm install --lts
fi
yarn config set registry https://registry.npm.taobao.org/
yarn

go mod vendor
