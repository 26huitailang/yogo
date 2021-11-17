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
}
