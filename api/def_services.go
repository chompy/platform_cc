package api

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

// serviceConfigs - list of service config type
var serviceConfigs = []ServiceConfig{
	MariadbService{},
	RedisService{},
}

// ServiceDef - defines a service
type ServiceDef struct {
	Name          string
	Type          string    `yaml:"type" json:"type"`
	Disk          int       `yaml:"disk" json:"disk"`
	Configuration yaml.Node `yaml:"configuration"`
}

// SetDefaults - set default values
func (d *ServiceDef) SetDefaults() {
	return
}

// Validate - validate service def
func (d ServiceDef) Validate() []error {
	o := make([]error, 0)
	if d.Type == "" {
		o = append(o, NewDefValidateError(
			"services[].type",
			"must be defined",
		))
	}
	serviceConfig := d.GetServiceConfig()
	if serviceConfig != nil {
		o = append(o, serviceConfig.Validate(&d)...)
	}
	return o
}

// GetServiceConfig - get configuration data for service
func (d ServiceDef) GetServiceConfig() ServiceConfig {
	for _, r := range serviceConfigs {
		if r.Check(&d) {
			return r
		}
	}
	return nil
}

// ParseServicesYaml - parse routes yaml
func ParseServicesYaml(d []byte) ([]*ServiceDef, error) {
	o := make(map[string]*ServiceDef)
	e := yaml.Unmarshal(d, &o)
	oo := make([]*ServiceDef, 0)
	for k := range o {
		o[k].SetDefaults()
		o[k].Name = k
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
