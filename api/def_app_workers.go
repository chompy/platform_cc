package api

// AppWorkerDef - defines a worker
type AppWorkerDef struct {
	Size          string                            `yaml:"size"`
	Disk          int                               `yaml:"disk"`
	Mounts        map[string]*AppMountDef           `yaml:"mounts"`
	Relationships map[string]string                 `yaml:"relationships"`
	Variables     map[string]map[string]interface{} `yaml:"variables"`
}

// SetDefaults - set default values
func (d *AppWorkerDef) SetDefaults() {
	for k := range d.Mounts {
		d.Mounts[k].SetDefaults()
	}
	if d.Size == "" {
		d.Size = "S"
	}
	if d.Disk < 256 {
		d.Disk = 256
	}
}

// Validate - validate AppWorkerDef
func (d AppWorkerDef) Validate() []error {
	// TODO
	return []error{}
}
