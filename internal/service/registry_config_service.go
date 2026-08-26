package service

import (
	"sort"
	"strings"
	"time"

	"serviceregistry/internal/model"
	"serviceregistry/pkg/idgen"
)

// SetConfig 新增或更新注册中心配置项。
func (s *Service) SetConfig(key, value, description string) (*model.RegistryConfig, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return nil, model.NewValidationError("key", "配置键不能为空")
	}
	if existing, err := s.store.GetRegistryConfigByKey(key); err == nil {
		existing.Value = value
		existing.Description = description
		existing.UpdatedAt = time.Now()
		if err := s.store.UpdateRegistryConfig(existing); err != nil {
			return nil, err
		}
		return existing, nil
	}
	c := &model.RegistryConfig{
		ID:          idgen.Hex(),
		Key:         key,
		Value:       value,
		Description: description,
		UpdatedAt:   time.Now(),
	}
	if err := s.store.CreateRegistryConfig(c); err != nil {
		return nil, err
	}
	return c, nil
}

// GetConfig 按键获取配置项。
func (s *Service) GetConfig(key string) (*model.RegistryConfig, error) {
	return s.store.GetRegistryConfigByKey(key)
}

// ListConfigs 返回全部配置项（按键排序）。
func (s *Service) ListConfigs() []*model.RegistryConfig {
	list := s.store.ListRegistryConfigs()
	sort.Slice(list, func(i, j int) bool { return list[i].Key < list[j].Key })
	return list
}

// DeleteConfig 按键删除配置项。
func (s *Service) DeleteConfig(key string) error {
	c, err := s.store.GetRegistryConfigByKey(key)
	if err != nil {
		return err
	}
	return s.store.DeleteRegistryConfig(c.ID)
}
