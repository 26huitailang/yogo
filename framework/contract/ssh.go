package contract

import (
	"fmt"

	"github.com/26huitailang/yogo/framework"
	"golang.org/x/crypto/ssh"
)

const SSHKey = "yogo:ssh"

type SSHOption func(container framework.Container, config *SSHConfig) error

type SSHService interface {
	GetClient(option ...SSHOption) (*ssh.Client, error)
}

type SSHConfig struct {
	Network string
	Host    string
	Port    string
	*ssh.ClientConfig
}

func (config *SSHConfig) UniqKey() string {
	return fmt.Sprintf("%v_%v_%v", config.Host, config.Port, config.User)
}
