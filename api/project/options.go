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

package project

// Option defines a project option.
type Option string

const (
	// OptionDomainSuffix sets the internal route domain suffix.
	OptionDomainSuffix Option = "domain_suffix"
)

// DefaultValue returns the default value of the option.
func (o Option) DefaultValue() string {
	switch o {
	case OptionDomainSuffix:
		{
			return "pcc.localtest.me"
		}
	}
	return ""
}

// Value returns the current value of the option with the default value if empty.
func (o Option) Value(opts map[Option]string) string {
	if opts[o] != "" {
		return opts[o]
	}
	return o.DefaultValue()
}

// ListOptions list all available project options.
func ListOptions() []Option {
	return []Option{
		OptionDomainSuffix,
	}
}

// ListOptionDescription returns a mapping of option name to its description.
func ListOptionDescription() map[Option]string {
	return map[Option]string{
		OptionDomainSuffix: "Domain name suffix for internal routes.",
	}
}
