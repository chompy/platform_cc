package project

// Option defines a project option.
type Option string

const (
	// OptionDNSIP set ip address that the internal pcc domain should point to.
	OptionDNSIP Option = "dns_ip"
)

// DefaultValue returns the default value of the option.
func (o Option) DefaultValue() string {
	switch o {
	case OptionDNSIP:
		{
			return "127.0.0.1"
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
		OptionDNSIP,
	}
}

// ListOptionDescription returns a mapping of option name to its description.
func ListOptionDescription() map[Option]string {
	return map[Option]string{
		OptionDNSIP: "IP address to use in PCC DNS server.",
	}
}
