package command

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/26huitailang/yogo/framework/cobra"
)

// build 相关命令
func initBuildCommand() *cobra.Command {
	buildCommand.AddCommand(buildFrontendCommand)
	buildCommand.AddCommand(buildBackendCommand)
	buildCommand.AddCommand(buildSelfCommand)
	buildCommand.AddCommand(buildAllCommand)
	return buildCommand
}

var buildCommand = &cobra.Command{
	Use:   "build",
	Short: "编译相关命令",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) == 0 {
			c.Help()
		}
		return nil
	},
}

// 打印前端的命令
var buildFrontendCommand = &cobra.Command{
	Use:   "frontend",
	Short: "使用yarn编译前端",
	RunE: func(c *cobra.Command, args []string) error {
		path, err := exec.LookPath("yarn")
		if err != nil {
			log.Fatalln("请安装yarn在你的PATH路径下")
		}

		cmd := exec.Command(path, "build")
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("=========== 前端编译失败 ===========")
			fmt.Println(string(out))
			fmt.Println("=========== 前端编译失败 ===========")
			return err
		}
		fmt.Print(string(out))
		fmt.Println("=========== 前端编译成功 ===========")
		return nil
	},
}

var buildSelfCommand = &cobra.Command{
	Use:   "self",
	Short: "编译yogo命令",
	RunE: func(c *cobra.Command, args []string) error {
		path, err := exec.LookPath("go")
		if err != nil {
			log.Fatalln("go: 请在Path路径中先安装go")
		}

		cmd := exec.Command(path, "build", "./")
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("go build error:")
			fmt.Println(string(out))
			fmt.Println("--------------")
			return err
		}
		fmt.Println("编译yogo成功")
		return nil
	},
}

// 打印后端的命令
var buildBackendCommand = &cobra.Command{
	Use:   "backend",
	Short: "使用go编译后端",
	RunE: func(c *cobra.Command, args []string) error {
		return buildSelfCommand.RunE(c, args)
	},
}

var buildAllCommand = &cobra.Command{
	Use:   "all",
	Short: "同时编译前端和后端",
	RunE: func(c *cobra.Command, args []string) error {
		err := buildFrontendCommand.RunE(c, args)
		if err != nil {
			return err
		}
		err = buildBackendCommand.RunE(c, args)
		if err != nil {
			return err
		}
		return nil
	},
}
