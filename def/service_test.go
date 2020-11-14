package def

import (
	"path"
	"strings"
	"testing"
)

func TestServices(t *testing.T) {
	p := path.Join("test_data", "sample1", ".platform", "services.yaml")
	services, e := ParseServicesYamlFile(p)
	if e != nil {
		t.Errorf("failed to parse services yaml, %s", e)
	}
	assertEqual(
		len(services),
		4,
		"unexpected number of services",
		t,
	)
	for _, service := range services {
		switch service.Name {
		case "mysqldb":
			{
				assertEqual(
					service.Type,
					"mysql:10.0",
					"unexpected services[].type",
					t,
				)
				assertEqual(
					service.Disk,
					512,
					"unexpected services[].disk",
					t,
				)
				break
			}
		case "rediscache":
			{
				assertEqual(
					service.Type,
					"redis:3.2",
					"unexpected services[].type",
					t,
				)
				assertEqual(
					service.Disk,
					0,
					"unexpected services[].disk",
					t,
				)
				break
			}
		case "redissession":
			{
				assertEqual(
					service.Type,
					"redis-persistent:3.2",
					"unexpected services[].type",
					t,
				)
				assertEqual(
					service.Disk,
					1024,
					"unexpected services[].disk",
					t,
				)
				break
			}
		case "solrsearch":
			{
				assertEqual(
					service.Type,
					"solr:6.6",
					"unexpected services[].type",
					t,
				)
				conf := service.Configuration["cores"].(map[string]interface{})["test"].(map[string]interface{})["conf_dir"].(string)
				assertEqual(
					strings.HasPrefix(conf, "H4sIAAAAAAAA"),
					true,
					"unexpected services[].configuration.cores[].conf_dir",
					t,
				)
				break
			}
		}

	}
}
