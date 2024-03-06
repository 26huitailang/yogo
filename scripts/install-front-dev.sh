#!/usr/bin/bash
curl -sL https://deb.nodesource.com/setup_20.x | bash -
apt-get install nodejs -y
node --version
npm install -g pnpm
pnpm --version
