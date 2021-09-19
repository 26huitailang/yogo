package main

import "github.com/26huitailang/yogo/framework"

func registerRouter(core *framework.Core) {
	core.Get("foo", FooControllerHandler)
}
