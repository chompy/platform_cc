package api

import (
	"io/ioutil"
	"log"
	"strings"

	"gopkg.in/yaml.v3"
)

// ServiceDef - defines a service
type ServiceDef struct {
	Name          string
	Type          string                  `yaml:"type" json:"type"`
	Disk          int                     `yaml:"disk" json:"disk"`
	Configuration ServiceConfigurationDef `yaml:"configuration"`
}

// SetDefaults - set default values
func (d *ServiceDef) SetDefaults() {
	if d.Configuration == nil {
		d.Configuration = make(map[string]interface{})
	}
	if d.Configuration["application_size"] == nil {
		d.Configuration["application_size"] = 0
	}
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
	return o
}

// GetTypeName - get type name
func (d ServiceDef) GetTypeName() string {
	return strings.Split(d.Type, ":")[0]
}

// GetEmptyRelationship - get empty relationship
func (d ServiceDef) GetEmptyRelationship() map[string]interface{} {
	return map[string]interface{}{
		"host":        "",
		"hostname":    "",
		"ip":          "",
		"port":        80,
		"path":        "",
		"scheme":      d.GetTypeName(),
		"fragment":    nil,
		"rel":         "",
		"host_mapped": false,
		"public":      false,
		"type":        d.Type,
		"service":     d.Name,
	}
}

// ParseServicesYaml - parse services yaml
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
