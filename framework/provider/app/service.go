package app

import (
	"errors"
	"flag"
	"path/filepath"

	"github.com/26huitailang/yogo/framework"
	"github.com/26huitailang/yogo/framework/util"
	"github.com/google/uuid"
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

	if baseFolder == "" {
		flag.StringVar(&baseFolder, "base", "", "base folder, default pwd")
		flag.Parse()
	}
	appId := uuid.New().String()
	configMap := map[string]string{}
	return &YogoApp{appId: appId, baseFolder: baseFolder, container: container, configMap: configMap}, nil
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

func (y YogoApp) HttpFolder() string {
	if val, ok := y.configMap["http_folder"]; ok {
		return val
	}
	return filepath.Join(y.AppFolder(), "http")
}

func (y YogoApp) MiddlewareFolder() string {
	if val, ok := y.configMap["middleware_folder"]; ok {
		return val
	}
	return filepath.Join(y.HttpFolder(), "middleware")
}

func (y YogoApp) ConsoleFolder() string {
	if val, ok := y.configMap["console_folder"]; ok {
		return val
	}
	return filepath.Join(y.BaseFolder(), "app", "console")
}

func (y YogoApp) CommandFolder() string {
	if val, ok := y.configMap["command_folder"]; ok {
		return val
	}
	return filepath.Join(y.ConsoleFolder(), "command")
}

func (y YogoApp) TestFolder() string {
	return filepath.Join(y.BaseFolder(), "test")
}

func (y YogoApp) BaseFolder() string {
	if y.baseFolder != "" {
		return y.baseFolder
	}

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

// AppFolder 代表app目录
func (app *YogoApp) AppFolder() string {
	if val, ok := app.configMap["app_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "app")
}

// DeployFolder 定义测试需要的信息
func (app *YogoApp) DeployFolder() string {
	if val, ok := app.configMap["deploy_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "deploy")
}
