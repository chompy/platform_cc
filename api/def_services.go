package api

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

// serviceList - list of service resolvers
var serviceResolvers = []ServiceResolver{
	MariadbService{},
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
	resolver := d.getServiceResolver()
	if resolver != nil {
		o = append(o, resolver.Validate(&d)...)
	}
	return o
}

// getServiceResolver - get service resolve for this service
func (d ServiceDef) getServiceResolver() ServiceResolver {
	for _, r := range serviceResolvers {
		if r.CheckType(d.Type) {
			return r
		}
	}
	return nil
}

// GetContainerConfig - get container configuration
func (d ServiceDef) GetContainerConfig() ServiceContainerDef {
	resolver := d.getServiceResolver()
	if resolver == nil {
		return nil
	}
	return resolver.GetContainerConfig(&d)
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
