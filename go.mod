module gitlab.com/contextualcode/platform_cc

go 1.15

replace gitlab.com/contextualcode/platform_cc/api/project => ./api/project

replace gitlab.com/contextualcode/platform_cc/api/docker => ./api/docker

replace gitlab.com/contextualcode/platform_cc/api/output => ./api/output

replace gitlab.com/contextualcode/platform_cc/api/def => ./api/def

replace gitlab.com/contextualcode/platform_cc/api/router => ./api/router

replace gitlab.com/contextualcode/platform_cc/cmd => ./cmd

replace gitlab.com/contextualcode/platform_cc/api/tests => ./api/tests

require (
	github.com/creack/pty v1.1.11 // indirect
	github.com/docker/cli v20.10.0-rc1+incompatible // indirect
	github.com/fatih/color v1.10.0 // indirect
	github.com/goccy/go-yaml v1.8.3 // indirect
	github.com/gopherjs/vecty v0.5.0 // indirect
	github.com/gorhill/cronexpr v0.0.0-20180427100037-88b0669f7d75
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/nsf/termbox-go v0.0.0-20200418040025-38ba6e5628f1
	github.com/olekukonko/tablewriter v0.0.4 // indirect
	github.com/robbiev/devdns v0.0.0-20141229153744-104b1f0d1b25 // indirect
	github.com/spf13/cobra v1.1.1
	github.com/ztrue/tracerr v0.3.0
	gitlab.com/contextualcode/platform_cc/api/def v0.0.1
	gitlab.com/contextualcode/platform_cc/api/docker v0.0.1
	gitlab.com/contextualcode/platform_cc/api/output v0.0.1
	gitlab.com/contextualcode/platform_cc/api/project v0.0.1
	gitlab.com/contextualcode/platform_cc/api/router v0.0.1
	gitlab.com/contextualcode/platform_cc/api/tests v0.0.1
	gitlab.com/contextualcode/platform_cc/cmd v0.0.1
	golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/yaml.v2 v2.3.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
)
