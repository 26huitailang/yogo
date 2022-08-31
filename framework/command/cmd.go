package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/26huitailang/yogo/framework/cobra"
	"github.com/26huitailang/yogo/framework/contract"
	"github.com/26huitailang/yogo/framework/util"
	"github.com/AlecAivazis/survey/v2"
	"github.com/jianfengye/collection"
	"github.com/pkg/errors"
)

// 初始化provider相关服务
func initCmdCommand() *cobra.Command {
	cmdCommand.AddCommand(cmdListCommand)
	cmdCommand.AddCommand(cmdCreateCommand)
	return cmdCommand
}

var cmdCommand = &cobra.Command{
	Use:   "cmd",
	Short: "命令提供相关命令",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) == 0 {
			c.Help()
		}
		return nil
	},
}

var cmdListCommand = &cobra.Command{
	Use:   "list",
	Short: "命令列表",
	RunE: func(c *cobra.Command, args []string) error {
		cmds := c.Root().Commands()
		info := [][]string{}
		var treeList func(cmds []*cobra.Command, level int)
		treeList = func(cmds []*cobra.Command, level int) {
			if len(cmds) == 0 {
				return
			}
			padding := strings.Repeat("    ", level-1) + " |---"
			for _, line := range cmds {
				info = append(info, []string{padding + line.Name(), line.Short})
				if len(line.Commands()) > 0 {
					level++
					treeList(line.Commands(), level)
				}
			}
		}
		for _, line := range cmds {
			info = append(info, []string{line.Name(), line.Short})
			treeList(line.Commands(), 1)
		}
		util.PrettyPrint(info)

		return nil
	},
}

// providerCreateCommand 创建一个新的服务，包括服务提供者，服务接口协议，服务实例
var cmdCreateCommand = &cobra.Command{
	Use:     "new",
	Aliases: []string{"create", "init"},
	Short:   "创建一个控制台命令",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		fmt.Println("开始创建控制台命令...")
		var name string
		var folder string
		{
			prompt := &survey.Input{
				Message: "请输入控制台命令名称:",
			}
			err := survey.AskOne(prompt, &name)
			if err != nil {
				return err
			}
		}
		{
			prompt := &survey.Input{
				Message: "请输入文件夹名称(默认: 同控制台命令):",
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

		pFolder := app.CommandFolder()
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

			t := template.Must(template.New("cmd").Funcs(funcs).Parse(cmdTmpl))
			if err := t.Execute(f, name); err != nil {
				return errors.Cause(err)
			}
		}
		fmt.Println("创建控制台命令工具陈工, 路径:", filepath.Join(pFolder, folder))
		fmt.Println("请记得开发完成后将命令行工具挂载到 console/kernel.go")
		return nil
	},
}

var cmdTmpl string = `package {{.}}

import (
	"fmt"

	"github.com/26huitailang/yogo/framework/cobra"
)

var {{.|title}}Command = &cobra.Command{
	Use:    "{{.}}",
	Short:  "{{.}}",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		fmt.Println(container)
		return nil
	},
}`
