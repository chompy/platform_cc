package api

import (
	"crypto/md5"
	"fmt"
	"io"
	"strings"
)

const mariadbPasswordSalt = "+AhX.v*H@lbVQ0)E4|d&xM#-zQR.68Y~Yy415Qp5C2,T"

// MariadbService - provides configuration for mariadb service
type MariadbService struct {
}

// Check - check if given service def matches
func (s MariadbService) Check(d *ServiceDef) bool {
	return strings.HasPrefix(d.Type, "mariadb") || strings.HasPrefix(d.Type, "mysql")
}

// Validate - validate service definition
func (s MariadbService) Validate(d *ServiceDef) []error {
	c, err := s.getConfiguration(d)
	if err != nil {
		return []error{err}
	}
	return c.Validate()
}

// getConfiguration - get configuration from yaml node
func (s MariadbService) getConfiguration(d *ServiceDef) (MariadbConfigurationDef, error) {
	o := MariadbConfigurationDef{}
	err := d.Configuration.Decode(&o)
	o.SetDefaults()
	return o, err
}

// getEndpointPassword - get password for an endpoint
func (s MariadbService) getEndpointPassword(user string) string {
	h := md5.New()
	io.WriteString(h, mariadbPasswordSalt)
	io.WriteString(h, user)
	io.WriteString(h, mariadbPasswordSalt)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// GetSetupCommand - get command to run to setup service
func (s MariadbService) GetSetupCommand(d *ServiceDef) ([]string, error) {
	// get config
	config, err := s.getConfiguration(d)
	if err != nil {
		return []string{}, err
	}
	// wait for mysqld
	cmd := `
sv mysql-standalone start
while ! mysqladmin ping --silent; do
    sleep 1
done
	`
	sqlQuery := ""
	// create schemas
	for _, s := range config.Schemas {
		sqlQuery += fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s CHARACTER SET UTF8mb4 COLLATE utf8mb4_bin; ", s)
	}
	// create users with privileges
	for user, end := range config.Endpoints {
		sqlQuery += fmt.Sprintf("CREATE USER '%s'@'%%'; ", user)
		for schema, priv := range end.Privileges {
			// grant priv + set password
			switch priv {
			case "admin":
				{
					sqlQuery += fmt.Sprintf(
						"GRANT ALL PRIVILEGES ON %s.* TO '%s'@'%%' IDENTIFIED BY '%s'; ",
						schema, user, s.getEndpointPassword(user),
					)
					break
				}
			case "ro":
				{
					sqlQuery += fmt.Sprintf(
						"GRANT SELECT ON %s.* TO '%s'@'%%' IDENTIFIED BY '%s'; ",
						schema, user, s.getEndpointPassword(user),
					)
					break
				}
			case "rw":
				{
					sqlQuery += fmt.Sprintf(
						"GRANT SELECT, INSERT, UPDATE, DELETE ON %s.* TO '%s'@'%%' IDENTIFIED BY '%s'; ",
						schema, user, s.getEndpointPassword(user),
					)
					break
				}
			}
		}
		sqlQuery += "FLUSH PRIVILEGES; "
		cmd += fmt.Sprintf(`mysql -h 127.0.0.1 -p$(cat /mnt/data/.mysql-password) -e "%s"`, sqlQuery)
	}
	return []string{"sh", "-c", cmd}, nil
}

// GetRelationship - get values for relationships variable
func (s MariadbService) GetRelationship(d *ServiceDef) ([]map[string]interface{}, error) {
	// get config
	config, err := s.getConfiguration(d)
	if err != nil {
		return []map[string]interface{}{}, err
	}
	out := make([]map[string]interface{}, 0)
	for user, end := range config.Endpoints {
		out = append(out, map[string]interface{}{
			"host":     "",
			"hostname": "",
			"ip":       "",
			"port":     3306,
			"path":     end.DefaultSchema,
			"username": user,
			"password": s.getEndpointPassword(user),
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
	return out, nil
}
