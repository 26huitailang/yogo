# yogo

- quasar: vue3 + vite
- gin: web framework
- cobra: cli

## todo

- [ ] casbin
- [ ] oauth2

## Project Tree

- [ ] confirm project menu usages
- [ ] confirm framework protocol

整个应用是围绕container构建的，符合 ServiceProvider定义就可以通过Bind方法注册到container中

最后，通过console作为程序入口启动对应的服务或者执行命令行

```tree
$ tree -I node_modules -I vendor -L 2
.
├── LICENSE
├── Makefile
├── README.md
├── app  # 构建真正的应用使用
│   ├── console  # 应用的命令行工具注册入口，同时包含framework的命令行
│   ├── http  # 应用http入口，注册应用路由和设置container，返回engine给framework YogoKernelProvider
│   └── provider  # 应用相关服务目录
├── bin
├── config  # 各种服务的配置，通过env.APP_ENV 申明使用那个配置
│   ├── development
│   ├── production
│   └── testing
├── deploy  # deploy指令存放文件的默认目录，可以通过配置app.deploy_folder修改
├── docker  # docker构建相关配置
│   └── Dockerfile
├── docs
│   ├── README.md
│   ├── guide
│   └── provider
├── framework  # framework核心文件
│   ├── cobra  # cobra源码 添加container相关操作
│   ├── command  # 框架cli
│   ├── container.go  # Container接口，YogoContain定义
│   ├── contract  # 服务定义
│   ├── gin  # gin源码 增加获取服务的一些方法，和定制request、response
│   ├── middleware  # 框架中间件
│   ├── provider  # provider服务实现
│   ├── provider.go  # provider接口定义
│   └── util  # 通用方法
├── go.mod
├── go.sum
├── index.html
├── main.go  # 程序入口
├── package.json
├── postcss.config.cjs
├── public
│   ├── favicon.ico
│   └── icons
├── quasar.config.js
├── scripts
│   └── install-dev-tools.sh
├── src
│   ├── App.vue
│   ├── assets
│   ├── boot
│   ├── components
│   ├── css
│   ├── env.d.ts
│   ├── i18n
│   ├── layouts
│   ├── pages
│   ├── quasar.d.ts
│   ├── router
│   ├── shims-vue.d.ts
│   └── stores
├── storage  # 运行时产生文件存储
│   ├── log
│   └── runtime
├── test
│   └── env.go
├── test.rest
├── tsconfig.json
└── yarn.lock

40 directories, 24 files
```

## dev

- devcontainer
