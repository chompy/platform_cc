package api

import (
	"strings"
)

// RedisService - provides configuration for redis service
type RedisService struct {
}

// Check - check if given service def matches
func (s RedisService) Check(d *ServiceDef) bool {
	return strings.HasPrefix(d.Type, "redis")
}

// Validate - validate service definition
func (s RedisService) Validate(d *ServiceDef) []error {
	return []error{}
}

// IsPersistent - check if this is persistent redis service
func (s RedisService) IsPersistent(d *ServiceDef) bool {
	return strings.HasPrefix(d.Type, "redis-persistent")
}

// GetSetupCommand - get command to run to setup service
func (s RedisService) GetSetupCommand(d *ServiceDef) ([]string, error) {
	return []string{}, nil
}

// GetRelationship - get values for relationships variable
func (s RedisService) GetRelationship(d *ServiceDef) ([]map[string]interface{}, error) {
	return []map[string]interface{}{
		map[string]interface{}{
			"host":        "",
			"hostname":    "",
			"ip":          "",
			"port":        6379,
			"path":        "",
			"scheme":      "redis",
			"fragment":    nil,
			"rel":         "redis",
			"host_mapped": false,
			"public":      false,
			"type":        "",
		},
	}, nil
}
