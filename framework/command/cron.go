package command

import (
	"fmt"
	"github.com/26huitailang/yogo/framework/cobra"
	"github.com/26huitailang/yogo/framework/contract"
	"github.com/26huitailang/yogo/framework/util"
	"github.com/erikdubbelboer/gspt"
	"github.com/sevlyar/go-daemon"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

var cronDaemon = false

func initCronCommand() *cobra.Command {
	cronStartCommand.Flags().BoolVarP(&cronDaemon, "daemon", "d", false, "start serve daemon")
	cronCommand.AddCommand(cronStartCommand)
	cronCommand.AddCommand(cronListCommand)
	return cronCommand
}

var cronCommand = &cobra.Command{
	Use:   "cron",
	Short: "定时任务相关命令",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) == 0 {
			c.Help()
		}
		return nil
	},
}

var cronStartCommand = &cobra.Command{
	Use:   "start",
	Short: "启动cron常驻进程",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		// 设置cron的日志地址和进程id地址
		pidFolder := appService.RuntimeFolder()
		serverPidFile := filepath.Join(pidFolder, "cron.pid")
		logFolder := appService.LogFolder()
		serverLogFile := filepath.Join(logFolder, "cron.log")
		currentFolder := appService.BaseFolder()

		if cronDaemon {
			cntxt := &daemon.Context{
				PidFileName: serverPidFile,
				PidFilePerm: 0664,
				LogFileName: serverLogFile,
				LogFilePerm: 0640,
				WorkDir:     currentFolder,
				// 设置所有设置文件的mask，默认为750
				Umask: 027,
				Args:  []string{"", "cron", "start", "--daemon=true"},
			}
			// 启动子进程，d不为空表示当前是父进程，d为空表示当前是子进程
			d, err := cntxt.Reborn()
			if err != nil {
				return err
			}
			if d != nil {
				// 父进程直接打印启动成功信息，不做任何操作
				fmt.Println("cron serve started, pid:", d.Pid)
				fmt.Println("log file", serverLogFile)
				return nil
			}

			// 子进程执行Cron.Run
			defer cntxt.Release()
			fmt.Println("daemon started")
			gspt.SetProcTitle("yogo cron")
			c.Root().Cron.Run()
			return nil
		}

		// 非daemon模式
		fmt.Println("start cron job")
		content := strconv.Itoa(os.Getpid())
		fmt.Println("[PID]", content)
		err := ioutil.WriteFile(serverPidFile, []byte(content), 0664)
		if err != nil {
			return err
		}

		gspt.SetProcTitle("yogo cron")
		c.Root().Cron.Run()
		return nil
	},
}

var cronListCommand = &cobra.Command{
	Use:   "list",
	Short: "列出所有的定时任务",
	RunE: func(c *cobra.Command, args []string) error {
		cronSpecs := c.Root().CronSpecs
		ps := [][]string{}
		for _, cronSpec := range cronSpecs {
			ps = append(ps, []string{cronSpec.Type, cronSpec.Spec, cronSpec.Cmd.Use, cronSpec.Cmd.Short, cronSpec.ServiceName})
		}
		util.PrettyPrint(ps)
		return nil
	},
}
