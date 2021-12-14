package command

import "github.com/26huitailang/yogo/framework/cobra"

func AddKernelCommands(root *cobra.Command) {
	root.AddCommand(initEnvCommand())
	root.AddCommand(DemoCommand)
	root.AddCommand(initAppCommand())
	root.AddCommand(initCronCommand())
}
