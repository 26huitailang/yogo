package command

import "github.com/26huitailang/yogo/framework/cobra"

// AddKernelCommands will add all command/* to root command
func AddKernelCommands(root *cobra.Command) {
	root.AddCommand(initNewCommand())
	root.AddCommand(initEnvCommand())
	root.AddCommand(initBuildCommand())
	root.AddCommand(initDevCommand())
	root.AddCommand(DemoCommand)
	root.AddCommand(initAppCommand())
	root.AddCommand(initCronCommand())
	root.AddCommand(initProviderCommand())
}
