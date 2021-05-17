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

import "fmt"

// Option defines a project option.
type Option string

const (
	// OptionDomainSuffix sets the internal route domain suffix.
	OptionDomainSuffix Option = "domain_suffix"
	// OptionMountStrategy defines the strategy of dealing with mounts.
	OptionMountStrategy Option = "mount_strategy"
)

const (
	// MountStrategyNone defines mount strategy where no action is taken.
	MountStrategyNone = "none"
	// MountStrategySymlink defines mount strategy where symlinks are used.
	MountStrategySymlink = "symlink"
	// MountStrategyVolume defines mount strategy where a container volume is used.
	MountStrategyVolume = "volume"
)

const (
	// OptionSourceLocal defines option source as local to current project.
	OptionSourceLocal = "local"
	// OptionSourceGlobal defines option source as global to the user.
	OptionSourceGlobal = "global"
	// OptionSourceNone defines option source as not set.
	OptionSourceNone = "unset"
)

// DefaultValue returns the default value of the option.
func (o Option) DefaultValue() string {
	switch o {
	case OptionDomainSuffix:
		{
			return "platform.cc"
		}
	case OptionMountStrategy:
		{
			return MountStrategyNone
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

// Validate returns error if given value is not a valid.
func (o Option) Validate(v string) error {
	switch o {
	case OptionMountStrategy:
		{
			if v == MountStrategyNone || v == MountStrategySymlink || v == MountStrategyVolume {
				return nil
			}
			return fmt.Errorf("mount strategy must be one of %s,%s,%s", MountStrategyNone, MountStrategySymlink, MountStrategyVolume)
		}
	}
	return nil

}

// ListOptions list all available project options.
func ListOptions() []Option {
	return []Option{
		OptionDomainSuffix,
		OptionMountStrategy,
	}
}

// ListOptionDescription returns a mapping of option name to its description.
func ListOptionDescription() map[Option]string {
	return map[Option]string{
		OptionDomainSuffix: "Domain name suffix for internal routes.",
		OptionMountStrategy: fmt.Sprintf(
			"Defines which mount strategy to use. (%s,%s,%s).",
			MountStrategyNone,
			MountStrategySymlink,
			MountStrategyVolume,
		),
	}
}

// GetOption returns the given option value globally or local to the project.
func (p *Project) GetOption(o Option) string {
	if p.Options[o] != "" {
		return p.Options[o]
	}
	gopt := p.globalConfig.Options[string(o)]
	if gopt != "" {
		return gopt
	}
	return o.DefaultValue()
}

// OptionSource returns the source of given option.
func (p *Project) OptionSource(o Option) string {
	if p.Options[o] != "" {
		return OptionSourceLocal
	}
	gopt := p.globalConfig.Options[string(o)]
	if gopt != "" {
		return OptionSourceGlobal
	}
	return OptionSourceNone
}
