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

// Flag defines flag that enable or disable features.
type Flag uint8

const (
	// EnableCron enables cron jobs.
	EnableCron Flag = 1 << iota
	// EnableWorkers enables workers.
	EnableWorkers
	// EnableServiceRoutes enables routes to services like Varnish.
	EnableServiceRoutes
	// EnablePHPOpcache enables PHP opcache.
	EnablePHPOpcache
)

// Add adds a flag.
func (f *Flag) Add(flag Flag) {
	*f = *f | flag
}

// Remove removes a flag.
func (f *Flag) Remove(flag Flag) {
	*f = *f &^ flag
}

// Has checks if flag is set.
func (f Flag) Has(flag Flag) bool {
	return f&flag != 0
}

// List returns a mapping of flag name to flag value.
func (f Flag) List() map[string]Flag {
	return map[string]Flag{
		"enable_cron":           EnableCron,
		"enable_workers":        EnableWorkers,
		"enable_service_routes": EnableServiceRoutes,
		"enable_php_opcache":    EnablePHPOpcache,
	}
}

// Descriptions returns a mapping of flag name to its description.
func (f Flag) Descriptions() map[string]string {
	return map[string]string{
		"enable_cron":           "Enables cron jobs.",
		"enable_workers":        "Enables workers.",
		"enable_service_routes": "Enable routes to services like Varnish.",
		"enable_php_opcache":    "Enables PHP Opcache.",
	}
}
