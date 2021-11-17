package contract

import "time"

const DistributedKey = "yogo:distributed"

// Distributed 分布式服务
type Distributed interface {
	// Select 分布式选择器，所有节点对某个服务进行抢占，只选择其中一个节点
	// ServiceName 服务名称
	// appId 当前的AppID
	// holdTime 抢占时间，单位秒
	// return selectAppID 分布式选择器最终选择的App, err 异常才返回，如果没有被选择，不返回
	Select(ServiceName string, appId string, holdTime time.Duration) (selectAppID string, err error)
}
