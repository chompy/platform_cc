package api

// AppRuntimeExtensionDef - defines an extension (PHP)
type AppRuntimeExtensionDef struct {
	Name          string            `yaml:"name" json:"name"`
	Configuration map[string]string `yaml:"configuration" json:"configuration,omitempty"`
}

// SetDefaults - set default values
func (d *AppRuntimeExtensionDef) SetDefaults() {
	return
}

// Validate - validate AppRuntimeExtensionDef
func (d AppRuntimeExtensionDef) Validate() []error {
	return []error{}
}

// UnmarshalYAML - implement Unmarshaler interface
func (d *AppRuntimeExtensionDef) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// unmarshal full extension
	data := make(map[string]interface{})
	e := unmarshal(&data)
	if e == nil {
		d.Name = data["name"].(string)
		d.Configuration = make(map[string]string)
		conf := data["configuration"].(map[string]interface{})
		for k, v := range conf {
			d.Configuration[k] = v.(string)
		}
		return nil
	}
	// unmarshal string extension name
	extName := ""
	e = unmarshal(&extName)
	if e != nil {
		return e
	}
	d.Name = extName
	d.Configuration = make(map[string]string)
	return nil
}
