package command

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/26huitailang/yogo/framework"
	"github.com/26huitailang/yogo/framework/cobra"
	"github.com/26huitailang/yogo/framework/contract"
	"github.com/26huitailang/yogo/framework/util"
	"github.com/erikdubbelboer/gspt"
	"github.com/sevlyar/go-daemon"
)

var appAddress = ""
var appDaemon = false

// initAppCommand 初始化app命令和其子命令
func initAppCommand() *cobra.Command {
	appStartCommand.Flags().BoolVarP(&appDaemon, "daemon", "d", false, "start app daemon")
	appStartCommand.Flags().StringVar(&appAddress, "address", "", "set app address, default: 8888")
	appCommand.AddCommand(appStartCommand)
	appCommand.AddCommand(appStateCommand)
	appCommand.AddCommand(appStopCommand)
	appCommand.AddCommand(appRestartCommand)
	appCommand.AddCommand(appVersionCommand)
	return appCommand
}

// AppCommand 是命令行参数第一级为app的命令，它没有实际功能，只是打印帮助文档
var appCommand = &cobra.Command{
	Use:   "app",
	Short: "业务应用控制命令",
	Long:  "业务应用控制命令，其包含业务启动，关闭，重启，查询等功能",
	RunE: func(c *cobra.Command, args []string) error {
		// 打印帮助文档
		c.Help()
		return nil
	},
}

func startAppServe(server *http.Server, c framework.Container) error {
	go func() {
		server.ListenAndServe()
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	closeWait := 5
	configService := c.MustMake(contract.ConfigKey).(contract.Config)
	if configService.IsExist("app.close_wait") {
		closeWait = configService.GetInt("app.close_wait")
	}
	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Duration(closeWait)*time.Second)
	defer cancel()

	if err := server.Shutdown(timeoutCtx); err != nil {
		return err
	}
	return nil
}

// appStartCommand 启动一个Web服务
var appStartCommand = &cobra.Command{
	Use:   "start",
	Short: "启动一个Web服务",
	RunE: func(c *cobra.Command, args []string) error {
		// 从Command中获取服务容器
		container := c.GetContainer()
		// 从服务容器中获取kernel的服务实例
		kernelService := container.MustMake(contract.KernelKey).(contract.Kernel)
		// 从kernel服务实例中获取引擎
		core := kernelService.HttpEngine()

		if appAddress == "" {
			envService := container.MustMake(contract.EnvKey).(contract.Env)
			if envService.Get("ADDRESS") != "" {
				appAddress = envService.Get("ADDRESS")
			} else {
				configService := container.MustMake(contract.ConfigKey).(contract.Config)
				if configService.IsExist("app.address") {
					appAddress = configService.GetString("app.address")
				} else {
					appAddress = ":8888"
				}
			}
		}
		// 创建一个Server服务
		server := &http.Server{
			Handler: core,
			Addr:    appAddress,
		}
		appService := container.MustMake(contract.AppKey).(contract.App)
		pidFolder := appService.RuntimeFolder()
		if !util.Exists(pidFolder) {
			if err := os.MkdirAll(pidFolder, os.ModePerm); err != nil {
				return err
			}
		}
		serverPidFile := filepath.Join(pidFolder, "app.pid")
		// check pid process exists
		if util.Exists(serverPidFile) {
			pidContent, err := os.ReadFile(serverPidFile)
			if err != nil {
				return err
			}
			if len(pidContent) > 0 {
				pid, err := strconv.Atoi(string(pidContent))
				if err != nil {
					return err
				}

				ok, _ := util.CheckProcessExistOptionalName(pid, "yogo")
				if ok {
					return fmt.Errorf("app is running: %d", pid)
				}
			}
		}
		logFolder := appService.LogFolder()
		if !util.Exists(logFolder) {
			if err := os.MkdirAll(logFolder, os.ModePerm); err != nil {
				return err
			}
		}
		serverLogFile := filepath.Join(logFolder, "app.log")
		currentFolder := util.GetExecDirectory()
		if appDaemon {
			cntxt := &daemon.Context{
				PidFileName: serverPidFile,
				PidFilePerm: 0664,
				LogFileName: serverLogFile,
				LogFilePerm: 0640,
				WorkDir:     currentFolder,
				Umask:       027,
				Args:        []string{"", "app", "start", "--daemon=true"},
			}
			d, err := cntxt.Reborn()
			if err != nil {
				return err
			}
			if d != nil {
				fmt.Println("app start successfully, pid:", d.Pid)
				fmt.Println("log file:", serverLogFile)
				return nil
			}
			defer cntxt.Release()
			fmt.Println("daemon started")
			if err := startAppServe(server, container); err != nil {
				fmt.Println(err)
			}
			return nil
		}

		// 非daemon模式
		content := strconv.Itoa(os.Getpid())
		fmt.Println("[PID]", content)
		err := os.WriteFile(serverPidFile, []byte(content), 0644)
		if err != nil {
			return err
		}
		gspt.SetProcTitle("yogo app")

		fmt.Println("app serve url:", appAddress)
		if err := startAppServe(server, container); err != nil {
			fmt.Println(err)
		}
		return nil
	},
}

var appStateCommand = &cobra.Command{
	Use:   "state",
	Short: "获取启动的app服务的pid",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)
		runtimeFolder := appService.RuntimeFolder()
		content, err := os.ReadFile(filepath.Join(runtimeFolder, "app.pid"))
		if err != nil {
			return err
		}
		if content != nil && len(content) > 0 {
			pid, err := strconv.Atoi(string(content))
			if err != nil {
				return err
			}
			if util.CheckProcessExist(pid) {
				fmt.Println("app is running, pid:", pid)
				return nil
			}
		}
		fmt.Println("no app service")
		return nil
	},
}
var appRestartCommand = &cobra.Command{
	Use:   "restart",
	Short: "应用状态",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)
		appPidFile := filepath.Join(appService.RuntimeFolder(), "app.pid")
		if !util.Exists(appPidFile) {
			fmt.Println("app.pid not found, start directly")
		} else {
			content, err := os.ReadFile(appPidFile)
			if err != nil {
				return err
			}
			if len(content) > 0 {
				pid, err := strconv.Atoi(string(content))
				if err != nil {
					return err
				}
				pExist, _ := util.CheckProcessExistOptionalName(pid, "yogo")
				if pExist {
					if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
						return err
					}
					closeWait := 5
					configService := container.MustMake(contract.ConfigKey).(contract.Config)
					if configService.IsExist("app.close_wait") {
						closeWait = configService.GetInt("app.close_wait")
					}
					for i := 0; i < closeWait*2; i++ {
						if !util.CheckProcessExist(pid) {
							break
						}
						time.Sleep(1 * time.Second)
					}
					if exist, _ := util.CheckProcessExistOptionalName(pid, "yogo"); exist {
						fmt.Println("stop app process failed, pid:", pid)
						return errors.New("stop app process failed")
					}

					if err := os.WriteFile(appPidFile, []byte(""), 0644); err != nil {
						return err
					}
					fmt.Println("stop app process successfully, pid:", pid)
				}
			}
		}
		appDaemon = false

		return appStartCommand.RunE(c, args)
	},
}

var appStopCommand = &cobra.Command{
	Use:   "stop",
	Short: "应用状态",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)
		appPidFile := filepath.Join(appService.RuntimeFolder(), "app.pid")
		content, err := os.ReadFile(appPidFile)
		if err != nil {
			return err
		}
		if len(content) > 0 {
			pid, err := strconv.Atoi(string(content))
			if err != nil {
				return err
			}
			if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
				return err
			}
			if err := os.WriteFile(appPidFile, []byte(""), 0644); err != nil {
				return err
			}
			fmt.Println("Stopped the app service, pid:", pid)
		}
		return nil
	},
}

var appVersionCommand = &cobra.Command{
	Use:   "version",
	Short: "yogo app version",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)
		fmt.Println(appService.Version())
		return nil
	},
}
