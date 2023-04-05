package command

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/26huitailang/yogo/framework/cobra"
	"github.com/26huitailang/yogo/framework/contract"
	"github.com/google/shlex"
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
		// 组装命令
		container := c.GetContainer()
		configService := container.MustMake(contract.ConfigKey).(contract.Config)
		appService := container.MustMake(contract.AppKey).(contract.App)
		logger := container.MustMake(contract.LogKey).(contract.Log)

		binFile := "yogo"
		if configService.GetString("deploy.backend.bin") != "" {
			binFile = configService.GetString("deploy.backend.bin")
		}
		useDocker := configService.GetBool("deploy.backend.use_docker")

		path := "/usr/local/go/bin/go"
		deployBinFile := filepath.Join("./", binFile)

		// deployBinFile := filepath.Join(deployFolder, binFile)
		if !useDocker {
			lpath, err := exec.LookPath("go")
			path = lpath
			if err != nil {
				log.Fatalln("yogo go: 请在Path路径中先安装go")
			}
		} else {
			// version := filepath.Base(deployFolder)
			deployFolderInDocker := "/code"
			deployBinFile = filepath.Join(deployFolderInDocker, binFile)
		}

		cmd := exec.Command(path, "build", "-o", deployBinFile, "./")
		cmd.Env = os.Environ()
		// 设置GOOS和GOARCH
		if configService.GetString("deploy.backend.goos") != "" {
			cmd.Env = append(cmd.Env, "GOOS="+configService.GetString("deploy.backend.goos"))
		}
		if configService.GetString("deploy.backend.goarch") != "" {
			cmd.Env = append(cmd.Env, "GOARCH="+configService.GetString("deploy.backend.goarch"))
		}
		if configService.GetString("deploy.backend.cgo") != "" {
			cmd.Env = append(cmd.Env, "CGO_ENABLED="+configService.GetString("deploy.backend.cgo"))
		}
		if configService.GetString("deploy.backend.cc") != "" {
			cmd.Env = append(cmd.Env, "CC="+configService.GetString("deploy.backend.cc"))
		}

		if useDocker {
			image := configService.GetString("deploy.backend.image")
			dcmd := fmt.Sprintf("run --rm -v %s:/code %s bash -c", filepath.Dir(appService.BaseFolder()), image)
			dcmds, _ := shlex.Split(dcmd)
			dcmds = append(dcmds, fmt.Sprintf("ls /code && cd /code && %s", cmd.String()))
			dockerPath, _ := exec.LookPath("docker")
			cmd = exec.Command(dockerPath, dcmds...)
		}

		// 执行命令
		ctx := context.Background()
		logger.Info(ctx, "开始执行命令:", map[string]interface{}{"command": cmd.String()})
		out, err := cmd.CombinedOutput()
		if err != nil {
			logger.Error(ctx, "go build err", map[string]interface{}{
				"err": err,
				"out": string(out),
			})
			return err
		}
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
