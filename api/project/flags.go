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

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"gitlab.com/contextualcode/platform_cc/api/output"
)

// FlagInt defines flag for enabled features. This is the old version that uses bitmasks.
type FlagInt uint8

const (
	// IntEnableCron enables cron jobs.
	IntEnableCron FlagInt = 1 << iota
	// IntEnableWorkers enables workers.
	IntEnableWorkers
	// IntEnableServiceRoutes enables routes to services like Varnish.
	IntEnableServiceRoutes
	// IntEnablePHPOpcache enables PHP opcache.
	IntEnablePHPOpcache
	// IntEnableMountVolume mounts Docker volume for Platform.sh mounts (NO LONGER USED).
	IntEnableMountVolume
	// IntEnableOSXNFSMounts uses NFS for mounts on OSX.
	IntEnableOSXNFSMounts
	// IntDisableYamlOverrides disables Platform.CC specific YAML override files.
	IntDisableYamlOverrides
	// IntDisableAutoCommit disables automatic commit of application container.
	IntDisableAutoCommit
)

// Add adds a flag.
func (f *FlagInt) Add(flag FlagInt) {
	*f = *f | flag
}

// Remove removes a flag.
func (f *FlagInt) Remove(flag FlagInt) {
	*f = *f &^ flag
}

// Has checks if flag is set.
func (f FlagInt) Has(flag FlagInt) bool {
	return f&flag != 0
}

// List returns a mapping of flag name to flag value.
func (f FlagInt) List() map[string]FlagInt {
	return map[string]FlagInt{
		EnableCron:           IntEnableCron,
		EnableWorkers:        IntEnableWorkers,
		EnableServiceRoutes:  IntEnableServiceRoutes,
		EnablePHPOpcache:     IntEnablePHPOpcache,
		EnableOSXNFSMounts:   IntEnableOSXNFSMounts,
		DisableYamlOverrides: IntDisableYamlOverrides,
		DisableAutoCommit:    IntDisableAutoCommit,
	}
}

// Flags defines flags that enable features.
type Flags map[string]int

const (
	// FlagUnset signifies that flag has not been set by user.
	FlagUnset = 0
	// FlagOff signifies that flag is set to off.
	FlagOff = 1
	// FlagOn signifies that flag is set to on.
	FlagOn = 2
)

const (
	// EnableCron enables cron jobs.
	EnableCron = "enable_cron"
	// EnableWorkers enables workers.
	EnableWorkers = "enable_workers"
	// EnableServiceRoutes enables routes to services like Varnish.
	EnableServiceRoutes = "enable_service_routes"
	// EnablePHPOpcache enables PHP opcache.
	EnablePHPOpcache = "enable_php_opcache"
	// EnableOSXNFSMounts uses NFS for mounts on OSX.
	EnableOSXNFSMounts = "enable_osx_nfs_mounts"
	// DisableYamlOverrides disables Platform.CC specific YAML override files.
	DisableYamlOverrides = "disable_yaml_overrides"
	// DisableAutoCommit disables automatic commit of application container.
	DisableAutoCommit = "disable_auto_commit"
	// DisableSharedGlobalVolume disables the shared global volume.
	DisableSharedGlobalVolume = "disable_shared_global_volume"
)

const (
	FlagSourceLocal  = "project"
	FlagSourceGlobal = "global"
	FlagSourceNone   = "unset"
)

// UnmarshalJSON implements Unmarshaler interface.
func (f *Flags) UnmarshalJSON(data []byte) error {
	// unmarshal as int (old flag set)
	// DEPRECATE
	dFlagInt := FlagInt(0)
	err := json.Unmarshal(data, &dFlagInt)
	if err == nil {
		(*f) = make(map[string]int)
		for name, flagIntVal := range dFlagInt.List() {
			(*f)[name] = FlagUnset
			if dFlagInt.Has(flagIntVal) {
				(*f)[name] = FlagOn
			}
		}
		return nil
	}
	// unmarshal as map
	dMap := make(map[string]int)
	err = json.Unmarshal(data, &dMap)
	if err != nil {
		return errors.WithStack(err)
	}
	*f = dMap
	return nil
}

// Descriptions returns a mapping of flag name to its description.
func (f Flags) Descriptions() map[string]string {
	return map[string]string{
		EnableCron:                "Enables cron jobs.",
		EnableWorkers:             "Enables workers.",
		EnableServiceRoutes:       "Enable routes to services like Varnish.",
		EnablePHPOpcache:          "Enables PHP Opcache.",
		EnableOSXNFSMounts:        "Enable NFS mounts on OSX.",
		DisableYamlOverrides:      "Disable Platform.CC specific YAML override files (.platform.app.pcc.yaml, services.pcc.yaml).",
		DisableAutoCommit:         "Disable auto commit of application containers on start.",
		DisableSharedGlobalVolume: "Disable the shared global volume.",
	}
}

// IsValidFlag returns true if given flag is valid.
func (f Flags) IsValidFlag(name string) bool {
	for iname := range f.Descriptions() {
		if name == iname {
			return true
		}
	}
	return false
}

// Set sets the value of the given flag.
func (f *Flags) Set(name string, value int) {
	if value < 0 {
		value = FlagUnset
	} else if value > 2 {
		value = FlagOn
	}
	if !f.IsValidFlag(name) {
		output.LogError(fmt.Errorf("tried to set unknown flag '%s' to %d", name, value))
		return
	}
	(*f)[name] = value
	output.LogDebug(fmt.Sprintf("Set flag '%s' to '%d.'", name, value), f)
}

// Get returns the value of the given flag.
func (f Flags) Get(name string) int {
	return f[name]
}

// IsOn returns true if given flag is on.
func (f Flags) IsOn(name string) bool {
	return f[name] == FlagOn
}

// IsOff returns true if given flag is off or unset.
func (f Flags) IsOff(name string) bool {
	return f[name] == FlagOff || f[name] == FlagUnset
}

// IsSet returns true if given flag has been set by the user.
func (f Flags) IsSet(name string) bool {
	return f[name] != FlagUnset
}

// IsUnset returns true if given flag has not been set by the user.
func (f Flags) IsUnset(name string) bool {
	return f[name] == FlagUnset
}

// GetValueName returns value of flag as a human readable string.
func (f Flags) GetValueName(name string) string {
	v := f.Get(name)
	if v < 0 {
		v = FlagUnset
	} else if v > 2 {
		v = FlagOn
	}
	switch v {
	case 1:
		{
			return "off"
		}
	case 2:
		{
			return "on"
		}
	default:
		{
			return "unset"
		}
	}
}

// HasFlag returns true if given flag is on globally or locally to the project.
func (p *Project) HasFlag(name string) bool {
	// check local project flag
	if p.Flags.IsSet(name) {
		return p.Flags.IsOn(name)
	}
	// check global flag
	for _, iname := range p.globalConfig.Flags {
		if name == iname {
			return true
		}
	}
	return false
}

// Source returns the source of the given flag.
func (p *Project) FlagSource(name string) string {
	// check local project flag
	if p.Flags.IsSet(name) {
		return FlagSourceLocal
	}
	// check global flag
	for _, iname := range p.globalConfig.Flags {
		if name == iname {
			return FlagSourceGlobal
		}
	}
	return FlagSourceNone
}
