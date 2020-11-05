package api

import (
	"crypto/md5"
	"fmt"
	"io"
	"strings"
)

const mariadbPasswordSalt = "j7pN8yT1zFgviBzvH2iwIKWbmGoQB7dm"

// MariadbContainer - defines configuration of container
type MariadbContainer struct {
	def           *ServiceDef
	configuration *MariadbConfigurationDef
}

// getEndpointPassword - get password for an endpoint
func (c MariadbContainer) getEndpointPassword(user string) string {
	h := md5.New()
	io.WriteString(h, mariadbPasswordSalt)
	io.WriteString(h, user)
	io.WriteString(h, mariadbPasswordSalt)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// GetImage - get image name
func (c MariadbContainer) GetImage() string {
	typeName := strings.Split(c.def.Type, ":")
	return fmt.Sprintf("%s%s-%s", platformShDockerImagePrefix, typeName[0], typeName[1])
}

// GetVolumes - get volumes used by service container
func (c MariadbContainer) GetVolumes() []string {
	return []string{"/var/lib/mysql"}
}

// GetStartCommand - get service container start command
func (c MariadbContainer) GetStartCommand() []string {
	cmd := `
mkdir -p /var/log/mysql
touch /var/log/mysql/mariadb-bin.index
chown -R mysql:mysql /var/log/mysql/
mkdir /var/run/mysqld || true
chown -R mysql:mysql /run/
usermod -s /bin/bash mysql
su mysql -c "/usr/sbin/mysqld"
	`
	return []string{"sh", "-c", cmd}
}

// GetPostStartCommand - get command to run after container is started
func (c MariadbContainer) GetPostStartCommand() []string {
	// wait for mysqld
	cmd := `
while ! mysqladmin ping --silent; do
    sleep 1
done
	`
	sqlQuery := ""
	// create schemas
	for _, s := range c.configuration.Schemas {
		sqlQuery += fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s CHARACTER SET UTF8mb4 COLLATE utf8mb4_bin; ", s)
	}
	// create users with privileges
	for user, end := range c.configuration.Endpoints {
		for schema, priv := range end.Privileges {
			// grant priv + set password
			switch priv {
			case "admin":
				{
					sqlQuery += fmt.Sprintf(
						"GRANT ALL PRIVILEGES ON %s.* TO '%s'@'%%' IDENTIFIED BY '%s'; ",
						schema, user, c.getEndpointPassword(user),
					)
					break
				}
			case "ro":
				{
					sqlQuery += fmt.Sprintf(
						"GRANT SELECT ON %s.* TO '%s'@'%%' IDENTIFIED BY '%s'; ",
						schema, user, c.getEndpointPassword(user),
					)
					break
				}
			case "rw":
				{
					sqlQuery += fmt.Sprintf(
						"GRANT SELECT, INSERT, UPDATE, DELETE ON %s.* TO '%s'@'%%' IDENTIFIED BY '%s'; ",
						schema, user, c.getEndpointPassword(user),
					)
					break
				}
			}
		}
		sqlQuery += "FLUSH PRIVILEGES; "
		cmd += fmt.Sprintf(`mysql -e "%s"`, sqlQuery)
	}
	return []string{"sh", "-c", cmd}
}

// GetRelationship - get values for relationships variable
func (c MariadbContainer) GetRelationship() []map[string]interface{} {
	out := make([]map[string]interface{}, 0)
	for user, end := range c.configuration.Endpoints {
		out = append(out, map[string]interface{}{
			"host":     "",
			"hostname": "",
			"ip":       "",
			"port":     3306,
			"path":     end.DefaultSchema,
			"username": user,
			"password": c.getEndpointPassword(user),
			"scheme":   "mysql",
			"fragment": nil,
			"query": map[string]interface{}{
				"is_master": true,
			},
			"rel":         user,
			"host_mapped": false,
			"public":      false,
			"type":        "",
		})
	}
	return out
}
