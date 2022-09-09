import { defaultTheme } from 'vuepress'
module.exports = {
  title: "yogo framework", // 设置网站标题
  description: "一个支持前后端开发的基于协议的框架", //描述
  dest: "./dist/", // 设置输出目录
  port: 2333, //端口
  base: "/",
  // head: [["link", { rel: "icon", href: "/assets/img/head.png" }]],
  theme: defaultTheme({
    //主题配置
    // logo: "/assets/img/head.png",
    repo: '26huitailang/yogo',
    docsDir: "docs",
    // 添加导航栏
    navbar: [
      { text: "主页", link: "/" }, // 导航条
      { text: "使用文档", link: "/guide/" },
      { text: "服务提供者", link: "/provider/" },
      {
        text: "github",
        // 这里是下拉列表展现形式。
        // children: [
        //   {
        //     text: "yogo",
        //     link: "https://github.com/26huitailang/yogo",
        //   },
        // ],
      },
    ],
    // 为以下路由添加侧边栏
    sidebar: {
      "/guide/": [
        {
          text: "指南",
          collapsible: false,
          children: [
            "introduce",
            "install",
            "build",
            "structure",
            "app",
            "env",
            "dev",
            "command",
            "cron",
            "middleware",
            "swagger",
            "provider",
            "todo",
          ],
        },
      ],
      "/provider/": [
        {
          text: "服务提供者",
          collapsible: false,
          children: [
            "app",
            "env",
            "config",
            "log",
          ],
        },
      ],
    },
  }),
};
