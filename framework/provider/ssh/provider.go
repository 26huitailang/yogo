package ssh

import (
	"github.com/26huitailang/yogo/framework"
	"github.com/26huitailang/yogo/framework/contract"
)

// SSHProvider 提供App的具体实现方法
type SSHProvider struct {
}

// Register 注册方法
func (h *SSHProvider) Register(container framework.Container) framework.NewInstance {
	return NewYogoSSH
}

// Boot 启动调用
func (h *SSHProvider) Boot(container framework.Container) error {
	return nil
}

// IsDefer 是否延迟初始化
func (h *SSHProvider) IsDefer() bool {
	return true
}

// Params 获取初始化参数
func (h *SSHProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container}
}

// Name 获取字符串凭证
func (h *SSHProvider) Name() string {
	return contract.SSHKey
}
