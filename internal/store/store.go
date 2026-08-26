// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"serviceregistry/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	// 服务
	CreateService(sv *model.Service) error
	GetService(id string) (*model.Service, error)
	ListServices() []*model.Service
	UpdateService(sv *model.Service) error
	DeleteService(id string) error

	// 实例
	CreateInstance(i *model.Instance) error
	GetInstance(id string) (*model.Instance, error)
	ListInstances() []*model.Instance
	UpdateInstance(i *model.Instance) error
	DeleteInstance(id string) error

	// 健康检查
	CreateHealthCheck(h *model.HealthCheck) error
	GetHealthCheck(id string) (*model.HealthCheck, error)
	GetHealthCheckByInstance(instanceID string) (*model.HealthCheck, error)
	ListHealthChecks() []*model.HealthCheck
	UpdateHealthCheck(h *model.HealthCheck) error
	DeleteHealthCheck(id string) error

	// 心跳
	CreateHeartbeat(h *model.Heartbeat) error
	ListHeartbeats() []*model.Heartbeat
	ListHeartbeatsByInstance(instanceID string) []*model.Heartbeat

	// 负载均衡策略
	CreateLoadBalanceStrategy(l *model.LoadBalanceStrategy) error
	GetLoadBalanceStrategy(id string) (*model.LoadBalanceStrategy, error)
	GetLoadBalanceByService(serviceID string) (*model.LoadBalanceStrategy, error)
	ListLoadBalanceStrategies() []*model.LoadBalanceStrategy
	UpdateLoadBalanceStrategy(l *model.LoadBalanceStrategy) error
	DeleteLoadBalanceStrategy(id string) error

	// 服务依赖
	CreateDependency(d *model.ServiceDependency) error
	GetDependency(id string) (*model.ServiceDependency, error)
	ListDependencies() []*model.ServiceDependency
	UpdateDependency(d *model.ServiceDependency) error
	DeleteDependency(id string) error

	// 注册事件
	CreateRegistryEvent(e *model.RegistryEvent) error
	ListRegistryEvents() []*model.RegistryEvent

	// 服务标签
	CreateServiceLabel(l *model.ServiceLabel) error
	GetServiceLabel(id string) (*model.ServiceLabel, error)
	ListServiceLabels() []*model.ServiceLabel
	UpdateServiceLabel(l *model.ServiceLabel) error
	DeleteServiceLabel(id string) error

	// 注册中心配置
	CreateRegistryConfig(c *model.RegistryConfig) error
	GetRegistryConfig(id string) (*model.RegistryConfig, error)
	GetRegistryConfigByKey(key string) (*model.RegistryConfig, error)
	ListRegistryConfigs() []*model.RegistryConfig
	UpdateRegistryConfig(c *model.RegistryConfig) error
	DeleteRegistryConfig(id string) error
}
