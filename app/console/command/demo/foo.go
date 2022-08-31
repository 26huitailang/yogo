package demo

import (
	"log"

	"github.com/26huitailang/yogo/framework/cobra"
)

func InitFooCommand() *cobra.Command {
	Foo3Command.AddCommand(Foo4Command)
	Foo2Command.AddCommand(Foo3Command)
	FooCommand.AddCommand(Foo2Command)
	return FooCommand
}

// FooCommand 代表Foo命令
var FooCommand = &cobra.Command{
	Use:     "foo",
	Short:   "foo的简要说明",
	Long:    "foo的长说明",
	Aliases: []string{"fo", "f"},
	Example: "foo命令的例子",
	RunE: func(c *cobra.Command, args []string) error {
		log.Println("execute foo command")
		return nil
	},
}

var Foo2Command = &cobra.Command{
	Use:     "foo2",
	Short:   "foo2的简要说明",
	Long:    "foo的长说明",
	Aliases: []string{"fo", "f"},
	Example: "foo命令的例子",
	RunE: func(c *cobra.Command, args []string) error {
		log.Println("execute foo command")
		return nil
	},
}
var Foo3Command = &cobra.Command{
	Use:     "foo3",
	Short:   "foo3的简要说明",
	Long:    "foo的长说明",
	Aliases: []string{"fo", "f"},
	Example: "foo命令的例子",
	RunE: func(c *cobra.Command, args []string) error {
		log.Println("execute foo command")
		return nil
	},
}
var Foo4Command = &cobra.Command{
	Use:     "foo4",
	Short:   "foo4的简要说明",
	Long:    "foo的长说明",
	Aliases: []string{"fo", "f"},
	Example: "foo命令的例子",
	RunE: func(c *cobra.Command, args []string) error {
		log.Println("execute foo command")
		return nil
	},
}
