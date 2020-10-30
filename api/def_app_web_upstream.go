package api

// AppWebUpstreamDef - defines how the front server will connect to the app
type AppWebUpstreamDef struct {
	SocketFamily string `yaml:"socket_family"`
	Protocol     string `yaml:"protocol"`
}

// SetDefaults - set default values
func (d *AppWebUpstreamDef) SetDefaults() {
	if d.SocketFamily == "" {
		d.SocketFamily = "tcp"
	}
	if d.Protocol == "" {
		d.Protocol = "fastcgi"
	}
}

// Validate - validate AppWebUpstreamDef
func (d AppWebUpstreamDef) Validate() []error {
	o := make([]error, 0)
	if d.SocketFamily != "tcp" && d.SocketFamily != "udp" {
		o = append(o, NewDefValidateError(
			"app.web.upstream.socket_family",
			"must be either tcp or udp",
		))
	}
	if d.Protocol != "http" && d.Protocol != "fastcgi" {
		o = append(o, NewDefValidateError(
			"app.web.upstream.protocol",
			"must be either http or fastcgi",
		))
	}
	return o
}
