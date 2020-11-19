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
