package app

import (
	"errors"
	"github.com/26huitailang/yogo/framework"
	"github.com/26huitailang/yogo/framework/util"
	"github.com/google/uuid"
	"path/filepath"
)

type YogoApp struct {
	appId      string
	container  framework.Container
	baseFolder string

	configMap map[string]string
}

func NewYogoApp(params ...interface{}) (interface{}, error) {
	if len(params) != 2 {
		return nil, errors.New("param error")
	}
	container := params[0].(framework.Container)
	baseFolder := params[1].(string)
	appId := uuid.New().String()
	return &YogoApp{appId: appId, baseFolder: baseFolder, container: container}, nil
}

func (y YogoApp) Version() string {
	return "0.0.1"
}

func (y *YogoApp) AppID() string {
	return y.appId
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

func (y YogoApp) TestFolder() string {
	return filepath.Join(y.BaseFolder(), "test")
}

func (y YogoApp) BaseFolder() string {
	if y.baseFolder != "" {
		return y.baseFolder
	}
	//var baseFolder string
	//flag.StringVar(&baseFolder, "base_folder", "", "base_folder param, default current path")
	//flag.Parsed()
	//if baseFolder != "" {
	//	return baseFolder
	//}

	return util.GetExecDirectory()
}

func (y YogoApp) StorageFolder() string {
	return filepath.Join(y.BaseFolder(), "storage")
}

func (y YogoApp) RuntimeFolder() string {
	return filepath.Join(y.StorageFolder(), "runtime")
}

func (y YogoApp) LogFolder() string {
	if val, ok := y.configMap["log_folder"]; ok {
		return val
	}
	return filepath.Join(y.LogFolder(), "log")
}

func (y *YogoApp) LoadAppConfig(kv map[string]string) {
	for key, val := range kv {
		y.configMap[key] = val
	}
}
