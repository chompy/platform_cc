/*
This file is part of Platform.CC.

Platform.CC is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Platform.CC is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Platform.CC.  If not, see <https://www.gnu.org/licenses/>.
*/

package def

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
	// TODO fix cron parser
	/*p := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	if _, e := p.Parse(d.Spec); e != nil {
		log.Println(d.Spec)
		o = append(o, NewValidateError(
			"app.cron[].spec",
			e.Error(),
		))
	}*/
	return o
}
