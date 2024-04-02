package auth

import (
	"fmt"

	"github.com/26huitailang/yogo/framework"
)

type AuthType string

const (
	AuthTypeOAuth2 AuthType = "oauth2"
)

type AuthService struct {
	container framework.Container
}

func NewAuthService(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	return &AuthService{container: container}, nil
}

// GetAuthenticator 根据需要获取对应的认证服务，如 oauth2、密码认证等
func (s *AuthService) GetAuthenticator(authType AuthType) (interface{}, error) {
	switch authType {
	case AuthTypeOAuth2:
		return NewOAuth2Service(s.container)
	default:
		return nil, fmt.Errorf("auth type not match: %s", authType)
	}
}

func NewOAuth2Service(params ...interface{}) (interface{}, error) {
	// todo: implement
	panic("implement me")
}
