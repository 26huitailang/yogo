package console

import (
	"github.com/26huitailang/yogo/app/console/command/demo"
	"github.com/26huitailang/yogo/framework"
	"github.com/26huitailang/yogo/framework/cobra"
	"github.com/26huitailang/yogo/framework/command"
)

func RunCommand(container framework.Container) error {
	var rootCmd = &cobra.Command{
		Use:   "yogo",
		Short: "yogo command",
		Long:  "yogo 框架命令行工具，使用这个命令可以执行框架自带的命令，使用--help查看详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.InitDefaultHelpCmd()
			return cmd.Help()
		},
		// todo: 尝试false是什么情况
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}
	rootCmd.SetContainer(container)
	command.AddKernelCommands(rootCmd)
	AddAppCommand(rootCmd)
	return rootCmd.Execute()
}

func AddAppCommand(rootCmd *cobra.Command) {
	rootCmd.AddCommand(demo.InitFoo())
	rootCmd.AddCronCommand("* * * * * *", demo.FooCommand)
}
