package api

import (
	"path"
	"testing"
)

func TestServices(t *testing.T) {
	p := path.Join("test_data", "sample1", "services.yaml")
	services, e := ParseServicesYamlFile(p)
	if e != nil {
		t.Errorf("failed to parse services yaml, %s", e)
	}
	assertEqual(
		services[0].Name,
		"mysqldb",
		"unexpected services[].name",
		t,
	)
	assertEqual(
		services[0].Type,
		"mysql:10.0",
		"unexpected services[].type",
		t,
	)
	assertEqual(
		services[0].Disk,
		512,
		"unexpected services[].disk",
		t,
	)
	assertEqual(
		services[0].Configuration.(ServicesMariadbConfigurationDef).Schemas[0],
		"main",
		"unexpected services[].configuration.schema[0]",
		t,
	)
	assertEqual(
		services[0].Configuration.(ServicesMariadbConfigurationDef).Endpoints["mysql"].DefaultSchema,
		"main",
		"unexpected services[].configuration.endpoints[].default_schema",
		t,
	)
	assertEqual(
		services[0].Configuration.(ServicesMariadbConfigurationDef).Endpoints["mysql"].Privileges["main"],
		"admin",
		"unexpected services[].configuration.endpoints[].privileges[]",
		t,
	)
}
