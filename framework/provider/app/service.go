package app

import (
	"errors"
	"flag"
	"github.com/26huitailang/yogo/framework"
	"github.com/26huitailang/yogo/framework/util"
	"path/filepath"
)

type YogoApp struct {
	container  framework.Container
	baseFolder string
}

func (y YogoApp) Version() string {
	return "0.0.1"
}

func (y YogoApp) ConfigFolder() string {
	return filepath.Join(y.BaseFolder(), "config")
}

func (y YogoApp) ProviderFolder() string {
	return filepath.Join(y.BaseFolder(), "app")
}

func (y YogoApp) MiddlewareFolder() string {
	return filepath.Join(y.BaseFolder(), "middleware")
}

func (y YogoApp) CommandFolder() string {
	return filepath.Join(y.BaseFolder(), "command")
}

func (y YogoApp) RuntimeFolder() string {
	return filepath.Join(y.BaseFolder(), "runtime")
}

func (y YogoApp) TestFolder() string {
	return filepath.Join(y.BaseFolder(), "test")
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

func NewYogoApp(params ...interface{}) (interface{}, error) {
	if len(params) != 2 {
		return nil, errors.New("param error")
	}
	container := params[0].(framework.Container)
	baseFolder := params[1].(string)
	return &YogoApp{baseFolder: baseFolder, container: container}, nil
}
