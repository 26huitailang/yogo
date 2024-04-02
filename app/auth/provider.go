package auth

import (
	"github.com/26huitailang/yogo/framework"
)

type AuthProvider struct {
	framework.ServiceProvider

	c framework.Container
}

func (sp *AuthProvider) Name() string {
	return AuthKey
}

func (sp *AuthProvider) Register(c framework.Container) framework.NewInstance {
	return NewAuthService
}

func (sp *AuthProvider) IsDefer() bool {
	return true
}

func (sp *AuthProvider) Params(c framework.Container) []interface{} {
	return []interface{}{c}
}

func (sp *AuthProvider) Boot(c framework.Container) error {
	return nil
}

