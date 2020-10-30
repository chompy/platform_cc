package api

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"gopkg.in/yaml.v2"
)

// ServiceDef - defines a service
type ServiceDef struct {
	Name          string
	Type          string      `yaml:"type"`
	Disk          int         `yaml:"disk"`
	Configuration interface{} `yaml:"configuration"`
}

// SetDefaults - set default values
func (d *ServiceDef) SetDefaults() {
	return
}

// Validate - validate ServiceDef
func (d ServiceDef) Validate() []error {
	o := make([]error, 0)
	if e := d.Configuration.(DefInterface).Validate(); len(e) > 0 {
		o = append(o, e...)
	}
	return o
}

// GetContainerImage - get container image for service
func (d ServiceDef) GetContainerImage() string {
	typeName := strings.Split(d.Type, ":")
	return fmt.Sprintf("%s%s-%s", platformShDockerImagePrefix, typeName[0], typeName[1])
}

// ParseServicesYaml - parse routes yaml
func ParseServicesYaml(d []byte) ([]*ServiceDef, error) {
	o := make(map[string]*ServiceDef)
	e := yaml.Unmarshal(d, &o)
	oo := make([]*ServiceDef, 0)
	for k := range o {
		o[k].SetDefaults()
		o[k].Name = k
		typeName := strings.Split(o[k].Type, ":")[0]
		switch typeName {
		case "mariadb", "mysql":
			{
				o[k].Configuration = mariadbConfigurationFromInterface(o[k].Configuration)
				break
			}
		}
		oo = append(oo, o[k])
	}
	return oo, e
}

// ParseServicesYamlFile - open services yaml file and parse it
func ParseServicesYamlFile(f string) ([]*ServiceDef, error) {
	log.Printf("Parse services at '%s.'", f)
	d, e := ioutil.ReadFile(f)
	if e != nil {
		return []*ServiceDef{}, e
	}
	return ParseServicesYaml(d)
}
