package model

import (
	"strings"
	"time"
)

// RegistryConfig 注册中心配置项（键值）。
type RegistryConfig struct {
	ID          string    `json:"id"`
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Validate 规范化并校验配置项字段。
func (c *RegistryConfig) Validate() error {
	c.Key = strings.TrimSpace(c.Key)
	c.Value = strings.TrimSpace(c.Value)
	c.Description = strings.TrimSpace(c.Description)
	if c.Key == "" {
		return NewValidationError("key", "配置键不能为空")
	}
	return nil
}
