package def

import (
	"github.com/gorhill/cronexpr"
)

// AppCron defines a cron job.
type AppCron struct {
	Spec    string `yaml:"spec" json:"spec"`
	Command string `yaml:"cmd" json:"cmd"`
}

// SetDefaults sets the default values.
func (d *AppCron) SetDefaults() {
	if d.Spec == "" {
		d.Spec = "* * * * *"
	}
}

// Validate checks for errors.
func (d AppCron) Validate(root *App) []error {
	o := make([]error, 0)
	if _, e := cronexpr.Parse(d.Spec); e != nil {
		o = append(o, NewDefValidateError(
			"app.cron[].spec",
			e.Error(),
		))
	}
	return o
}
