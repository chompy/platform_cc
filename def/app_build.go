package def

// AppBuild defines what happens when building the app.
type AppBuild struct {
	Flavor string `yaml:"flavor"`
}

// SetDefaults sets the default values.
func (d *AppBuild) SetDefaults() {
	if d.Flavor == "" {
		d.Flavor = "none"
	}
}

// Validate checks for errors.
func (d AppBuild) Validate(root *App) []error {
	o := make([]error, 0)
	switch root.GetTypeName() {
	case "php":
		{
			if err := validateMustContainOne(
				[]string{"composer", "drupal", "symfony", "default", "none", "update"},
				d.Flavor,
				"app.build.flavor",
			); err != nil {
				o = append(o, err)
			}
			break
		}
	}
	return o
}
