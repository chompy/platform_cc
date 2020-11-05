package api

import (
	"path"
	"strings"
	"testing"
)

func TestServices(t *testing.T) {
	p := path.Join("test_data", "sample1", "services.yaml")
	services, e := ParseServicesYamlFile(p)
	if e != nil {
		t.Errorf("failed to parse services yaml, %s", e)
	}
	assertEqual(
		len(services),
		3,
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
				assertEqual(
					service.GetContainerConfig().GetImage(),
					platformShDockerImagePrefix+"mariadb-10.0",
					"unexpected services[].GetContainerConfig.GetImage",
					t,
				)
				assertEqual(
					service.GetContainerConfig().GetVolumes()[0],
					"/var/lib/mysql",
					"unexpected service container config GetVolumes",
					t,
				)
				relationships := service.GetContainerConfig().GetRelationship()
				assertEqual(
					relationships[0]["username"],
					"mysql",
					"unexpected service container config GetRelationship[].username",
					t,
				)
				assertEqual(
					strings.Contains(
						strings.Join(service.GetContainerConfig().GetStartCommand(), " "),
						"mysqld",
					),
					true,
					"expected service container config GetStartCommand to contain string 'mysqld'",
					t,
				)
				assertEqual(
					strings.Contains(
						strings.Join(service.GetContainerConfig().GetPostStartCommand(), " "),
						"CREATE SCHEMA IF NOT EXISTS main CHARACTER SET UTF8mb4 COLLATE utf8mb4_bin",
					),
					true,
					"expected service container config GetPostStartCommand to contain schema create query",
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
		}

	}
}
