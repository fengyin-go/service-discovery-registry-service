package store

import (
	"sync"

	"serviceregistry/internal/model"
)

// MemoryStore 基于内存的 Store 实现，线程安全。
type MemoryStore struct {
	mu           sync.RWMutex
	services     map[string]*model.Service
	instances    map[string]*model.Instance
	healthChecks map[string]*model.HealthCheck
	heartbeats   map[string]*model.Heartbeat
	loadBalances map[string]*model.LoadBalanceStrategy
	dependencies map[string]*model.ServiceDependency
	events       map[string]*model.RegistryEvent
	labels       map[string]*model.ServiceLabel
	configs      map[string]*model.RegistryConfig
}

// NewMemoryStore 构造空的内存存储。
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		services:     make(map[string]*model.Service),
		instances:    make(map[string]*model.Instance),
		healthChecks: make(map[string]*model.HealthCheck),
		heartbeats:   make(map[string]*model.Heartbeat),
		loadBalances: make(map[string]*model.LoadBalanceStrategy),
		dependencies: make(map[string]*model.ServiceDependency),
		events:       make(map[string]*model.RegistryEvent),
		labels:       make(map[string]*model.ServiceLabel),
		configs:      make(map[string]*model.RegistryConfig),
	}
}

var _ Store = (*MemoryStore)(nil)
