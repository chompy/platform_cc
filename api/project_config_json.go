package api

import (
	"encoding/json"
)

// BuildConfigJSON - make config.json for container runtime
func (p *Project) BuildConfigJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"primary_ip": "127.0.0.1",
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
		"domainname":   "rqetxnydnu6llao4tkfrmkzesu.app.service._.us-2.platformsh.site",
		"host_ip":      "127.0.0.1",
		"applications": p.Apps,
		"configuration": map[string]interface{}{
			"access": map[string]interface{}{
				"ssh": []string{},
			},
			"privileged_digest": "44136fa355b3678a1146ad16f7e8649e94fb4fc21fe77e8310c060f61caaff8a",
			"environment_info": map[string]interface{}{
				"is_production": false,
				"machine_name":  "pcc-1",
				"name":          "pcc",
				"reference":     "refs/heads/pcc",
				"is_main":       false,
			},
			"project_info": map[string]interface{}{
				"name": p.ID,
				"settings": map[string]interface{}{
					"systemd":          false,
					"variables_prefix": "PLATFORM_",
					"crons_in_git":     false,
					"product_code":     "platformsh",
					"product_name":     "Platform.sh",
					"enforce_mfa":      false,
					"bot_email":        "bot@platform.sh",
				},
			},
			"privileged": map[string]interface{}{},
		},
		"info": map[string]interface{}{
			"mail_relay_host":    "",
			"mail_relay_host_v2": "127.0.0.1",
			"limits": map[string]interface{}{
				"disk":   p.Apps[0].Disk,
				"cpu":    1.0,
				"memory": 1024,
			},
			"external ip": "127.0.0.1",
		},
		"name":       p.ID,
		"service":    "app",
		"cluster":    "-",
		"region":     "us-2.platform.sh",
		"hostname":   "app.0",
		"instance":   p.ID,
		"nameserver": "1.1.1.1",
	})

}
