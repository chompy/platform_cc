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

// AppWebUpstream defines how the front server will connect to the app.
type AppWebUpstream struct {
	SocketFamily string `yaml:"socket_family" json:"socket_family,omitempty"`
	Protocol     string `yaml:"protocol" json:"protocol,omitempty"`
}

// SetDefaults sets the default values.
func (d *AppWebUpstream) SetDefaults() {
	// no defaults, the PSH container scripts can figure out the best configuration when omitted
}

// Validate checks for errors.
func (d AppWebUpstream) Validate(root *App) []error {
	o := make([]error, 0)
	if err := validateMustContainOne(
		[]string{"", "tcp", "udp", "unix"},
		d.SocketFamily,
		"app.web.upstream.socket_family",
	); err != nil {
		o = append(o, err)
	}
	if err := validateMustContainOne(
		[]string{"", "http", "fastcgi"},
		d.Protocol,
		"app.web.upstream.protocol",
	); err != nil {
		o = append(o, err)
	}
	return o
}
