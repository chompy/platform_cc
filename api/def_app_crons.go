package api

import (
	"github.com/gorhill/cronexpr"
)

// AppCronDef - defines a cron job
type AppCronDef struct {
	Spec    string `yaml:"spec" json:"spec"`
	Command string `yaml:"cmd" json:"cmd"`
}

// SetDefaults - set default values
func (d *AppCronDef) SetDefaults() {
	if d.Spec == "" {
		d.Spec = "* * * * *"
	}
}

// Validate - validate AppCronDef
func (d AppCronDef) Validate() []error {
	o := make([]error, 0)
	if _, e := cronexpr.Parse(d.Spec); e != nil {
		o = append(o, NewDefValidateError(
			"app.cron[].spec",
			e.Error(),
		))
	}
	return o
}
