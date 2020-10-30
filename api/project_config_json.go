package api

import (
	"encoding/json"
)

// BuildConfigJSON - make config.json for container runtime
func (p *Project) BuildConfigJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"primary_ip": "-",
		"features": []string{
			"extensible_cluster",
			"backup_unnamed",
			"service_relationships_extended",
			"service_instances",
			"hot_backup",
			"mounts_ng",
			"notify",
			"git_state_api",
			"instance_notify",
			"never_reapable_backup",
			"service_slug",
			"notify_cron",
			"prepare_state_deployment",
			"outbound_restrictions",
			"privileged_configuration",
			"configurable_build_resources",
			"stop_the_world_backup",
			"export_import_backup",
			"supports_state",
		},
		"domainname":   "-",
		"host_ip":      "-",
		"applications": p.Apps,
		"configuration": map[string]interface{}{
			"access": map[string]interface{}{
				"ssh": []string{},
			},
			"privileged_digest": "-",
			"environment_info": map[string]interface{}{
				"is_production": false,
				"machine_name":  "-",
				"name":          "pcc",
				"reference":     "-",
				"is_main":       false,
			},
			"project_info": map[string]interface{}{
				"name": p.ID,
				"settings": map[string]interface{}{
					"systemd":          false,
					"variables_prefix": "PLATFORM_",
					"crons_in_git":     false,
					"product_code":     "platformcc",
					"product_name":     "Platform.cc",
					"enforce_mfa":      false,
					"bot_email":        "bot@pcc.ccplatform.net",
				},
			},
			"privileged": map[string]interface{}{},
		},
		"info": map[string]interface{}{
			"mail_relay_host":    "",
			"mail_relay_host_v2": "-",
			"limits": map[string]interface{}{
				"disk":   99999,
				"cpu":    1.0,
				"memory": 1024,
			},
			"external ip": "-",
		},
		"name":       p.ID,
		"service":    "app",
		"cluster":    "-",
		"region":     "-",
		"hostname":   "-",
		"instance":   "-",
		"nameserver": "1.1.1.1",
	})
}
