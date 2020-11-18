package def

// AppWorker defines a worker.
type AppWorker struct {
	Size          string                            `yaml:"size"`
	Disk          int                               `yaml:"disk"`
	Mounts        map[string]*AppMount              `yaml:"mounts"`
	Relationships map[string]string                 `yaml:"relationships"`
	Variables     map[string]map[string]interface{} `yaml:"variables"`
}

// SetDefaults sets the default values.
func (d *AppWorker) SetDefaults() {
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

// Validate checks for errors.
func (d AppWorker) Validate() []error {
	// TODO
	return []error{}
}
