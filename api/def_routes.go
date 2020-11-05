package api

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

// RouteDef - define a route
type RouteDef struct {
	Path       string
	Type       string            `yaml:"type"`
	Upstream   string            `yaml:"upstream"`
	To         string            `yaml:"to"`
	ID         string            `yaml:"id"`
	Attributes map[string]string `yaml:"attributes"`
	Cache      RouteCacheDef     `yaml:"cache"`
	Redirects  RouteRedirectsDef `yaml:"redirects"`
	SSI        RoutesSsiDef      `yaml:"ssi"`
}

// SetDefaults - set default values
func (d *RouteDef) SetDefaults() {
	if d.Type == "" {
		d.Type = "upstream"
	}
	d.Cache.SetDefaults()
	d.Redirects.SetDefaults()
	d.SSI.SetDefaults()
}

// Validate - validate RouteDef
func (d RouteDef) Validate() []error {
	o := make([]error, 0)
	if e := d.Cache.Validate(); e != nil {
		o = append(o, e...)
	}
	if e := d.Redirects.Validate(); e != nil {
		o = append(o, e...)
	}
	if e := d.SSI.Validate(); e != nil {
		o = append(o, e...)
	}
	return o
}

// ParseRoutesYaml - parse routes yaml
func ParseRoutesYaml(d []byte) ([]*RouteDef, error) {
	o := make(map[string]*RouteDef)
	e := yaml.Unmarshal(d, &o)
	oo := make([]*RouteDef, 0)
	for k := range o {
		o[k].SetDefaults()
		o[k].Path = k
		oo = append(oo, o[k])
	}
	return oo, e
}

// ParseRoutesYamlFile - open routes yaml file and parse it
func ParseRoutesYamlFile(f string) ([]*RouteDef, error) {
	log.Printf("Parse routes at '%s.'", f)
	d, e := ioutil.ReadFile(f)
	if e != nil {
		return []*RouteDef{}, e
	}
	return ParseRoutesYaml(d)
}
