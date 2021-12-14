package contract

const (
	EnvProduction  = "production"
	EnvTesting     = "testing"
	EnvDevelopment = "development"
	EnvKey         = "yogo:env"
)

type Env interface {
	AppEnv() string
	IsExist(key string) bool
	Get(string) string
	All() map[string]string
}
