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
	}
}

// Descriptions returns a mapping of flag name to its description.
func (f Flag) Descriptions() map[string]string {
	return map[string]string{
		"enable_cron":           "Enables cron jobs.",
		"enable_workers":        "Enables workers.",
		"enable_service_routes": "Enable routes to services like Varnish.",
	}
}
