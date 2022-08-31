package command

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/26huitailang/yogo/framework/cobra"
	"github.com/26huitailang/yogo/framework/contract"
	"github.com/26huitailang/yogo/framework/util"
	"github.com/AlecAivazis/survey/v2"
	"github.com/go-git/go-git/v5"
	"github.com/jianfengye/collection"
	"github.com/pkg/errors"
)

// 初始化provider相关服务
func initMiddlewareCommand() *cobra.Command {
	middlewareCommand.AddCommand(middlewareListCommand)
	middlewareCommand.AddCommand(middlewareCreateCommand)
	middlewareCommand.AddCommand(middlewareMigrateCommand)
	return middlewareCommand
}

var middlewareCommand = &cobra.Command{
	Use:   "middleware",
	Short: "中间件提供相关命令",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) == 0 {
			c.Help()
		}
		return nil
	},
}

var middlewareListCommand = &cobra.Command{
	Use:   "list",
	Short: "显示所有中间件",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		middlewarePath := path.Join(appService.BaseFolder(), "app", "http", "middleware")

		// 读取文件夹
		files, err := ioutil.ReadDir(middlewarePath)
		if err != nil {
			return err
		}
		for _, f := range files {
			if f.IsDir() {
				fmt.Println(f.Name())
			}
		}
		return nil
	},
}

var middlewareMigrateCommand = &cobra.Command{
	Use:   "migrate",
	Short: "迁移gin-contrib中间件, 迁移地址: https://github.com/gin-contrib/[middleware].git",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		// 获取中间件名称
		var name string
		{
			prompt := &survey.Input{
				Message: "请输入中间件名称: ",
			}
			err := survey.AskOne(prompt, &name)
			if err != nil {
				return err
			}
		}
		// 使用go-git 将对应gin-contrib项目clone 到 app/http/middleware目录
		app := container.MustMake(contract.AppKey).(contract.App)
		middlewarePath := app.MiddlewareFolder()
		url := fmt.Sprintf("https://github.com/gin-contrib/%s.git", name)
		fmt.Println("下载中间件 gint-contrib:")
		fmt.Println(url)
		_, err := git.PlainClone(path.Join(middlewarePath, name), false, &git.CloneOptions{
			URL:      url,
			Progress: os.Stdout,
		})
		if err != nil {
			return err
		}
		// 删除不必要的文件 go.mod go.sum .git/
		fileToRemove := map[string]bool{
			"go.mod": false,
			"go.sum": false,
			".git":   true,
		}
		repoFolder := path.Join(middlewarePath, name)
		for k, isFolder := range fileToRemove {
			path2rm := path.Join(repoFolder, k)
			fmt.Println("remove " + path2rm)
			if isFolder {
				os.RemoveAll(path2rm)
			} else {
				os.Remove(path2rm)
			}
		}
		// 替换关键字
		filepath.Walk(repoFolder, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			if filepath.Ext(path) != ".go" {
				return nil
			}
			c, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			isContain := bytes.Contains(c, []byte("github.com/gin-gonic/gin"))
			if isContain {
				fmt.Println("update: " + path)
				c = bytes.ReplaceAll(c, []byte("github.com/gin-gonic/gin"), []byte("github.com/26huitailang/yogo/framework/gin"))
				err = ioutil.WriteFile(path, c, 0644)
				if err != nil {
					return err
				}
			}
			return nil
		})
		return nil
	},
}

var middlewareCreateCommand = &cobra.Command{
	Use:     "new",
	Aliases: []string{"create", "init"},
	Short:   "创建一个中间件",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		fmt.Println("开始创建中间件...")
		var name string
		var folder string
		{
			prompt := &survey.Input{
				Message: "请输入中间件名称:",
			}
			err := survey.AskOne(prompt, &name)
			if err != nil {
				return err
			}
		}
		{
			prompt := &survey.Input{
				Message: "请输入文件夹名称(默认: 同中间件名称):",
			}
			err := survey.AskOne(prompt, &folder)
			if err != nil {
				return err
			}
		}

		if folder == "" {
			folder = name
		}
		// 检查文件是否存在
		app := container.MustMake(contract.AppKey).(contract.App)

		pFolder := app.MiddlewareFolder()
		subFolders, err := util.SubDir(pFolder)
		if err != nil {
			return err
		}
		subColl := collection.NewStrCollection(subFolders)
		if subColl.Contains(folder) {
			fmt.Println("目录名称已经存在")
			return nil
		}

		// 开始创建文件
		if err := os.Mkdir(filepath.Join(pFolder, folder), 0700); err != nil {
			return err
		}
		// 创建title这个模版方法
		funcs := template.FuncMap{"title": strings.Title}
		{
			file := filepath.Join(pFolder, folder, name+".go")
			f, err := os.Create(file)
			if err != nil {
				return errors.Cause(err)
			}

			t := template.Must(template.New("cmd").Funcs(funcs).Parse(middlewareTmp))
			if err := t.Execute(f, name); err != nil {
				return errors.Cause(err)
			}
		}
		fmt.Println("创建中间件成功, 路径:", filepath.Join(pFolder, folder))
		return nil
	},
}

var middlewareTmp string = `package {{.}}

import "github.com/26huitailang/yogo/framework/gin"

// {{.|title}}Middleware 代表中间件函数
func {{.|title}}Middleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Next()
	}
}

`
