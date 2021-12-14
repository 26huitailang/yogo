package contract

const AppKey = "yogo:app"

type App interface {
	AppID() string
	Version() string
	BaseFolder() string
	ConfigFolder() string
	LogFolder() string
	ProviderFolder() string
	MiddlewareFolder() string
	CommandFolder() string
	RuntimeFolder() string
	TestFolder() string

	// LoadAppConfig 加载新的AppConfig, key为对应函数转为小写下划线，比如ConfigFolder => config_folder
	LoadAppConfig(kv map[string]string)
}
