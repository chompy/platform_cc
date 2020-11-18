package def

// AppWebUpstream defines how the front server will connect to the app.
type AppWebUpstream struct {
	SocketFamily string `yaml:"socket_family"`
	Protocol     string `yaml:"protocol"`
}

// SetDefaults sets the default values.
func (d *AppWebUpstream) SetDefaults() {
	if d.SocketFamily == "" {
		d.SocketFamily = "tcp"
	}
	if d.Protocol == "" {
		d.Protocol = "fastcgi"
	}
}

// Validate checks for errors.
func (d AppWebUpstream) Validate(root *App) []error {
	o := make([]error, 0)
	if err := validateMustContainOne(
		[]string{"tcp", "udp"},
		d.SocketFamily,
		"app.web.upstream.socket_family",
	); err != nil {
		o = append(o, err)
	}
	if err := validateMustContainOne(
		[]string{"http", "fastcgi"},
		d.Protocol,
		"app.web.upstream.protocol",
	); err != nil {
		o = append(o, err)
	}
	return o
}
