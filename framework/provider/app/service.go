package app

import (
	"flag"
	"github.com/26huitailang/yogo/framework"
	"github.com/26huitailang/yogo/framework/util"
	"path/filepath"
)

type YogoApp struct {
	container  framework.Container
	baseFolder string
}

func (y YogoApp) BaseFolder() string {
	if y.baseFolder != "" {
		return y.baseFolder
	}
	var baseFolder string
	flag.StringVar(&baseFolder, "base_folder", "", "base_folder param, default current path")
	flag.Parsed()
	if baseFolder != "" {
		return baseFolder
	}

	return util.GetExecDirectory()
}

func (y YogoApp) StorageFolder() string {
	return filepath.Join(y.BaseFolder(), "storage")
}

func (y YogoApp) LogFolder() string {
	return filepath.Join(y.BaseFolder(), "log")
}
