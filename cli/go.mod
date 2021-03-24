module gitlab.com/contextualcode/platform_cc/cli

go 1.15

replace gitlab.com/contextualcode/platform_cc/api/project => ../api/project

replace gitlab.com/contextualcode/platform_cc/api/output => ../api/output

replace gitlab.com/contextualcode/platform_cc/api/def => ../api/def

replace gitlab.com/contextualcode/platform_cc/api/router => ../api/router

replace gitlab.com/contextualcode/platform_cc/api/container => ../api/container

require (
	github.com/olekukonko/tablewriter v0.0.4
	github.com/spf13/cobra v1.1.3
	github.com/ztrue/tracerr v0.3.0
	gitlab.com/contextualcode/platform_cc/api/container v0.0.1
	gitlab.com/contextualcode/platform_cc/api/def v0.0.1
	gitlab.com/contextualcode/platform_cc/api/output v0.0.1
	gitlab.com/contextualcode/platform_cc/api/project v0.0.1
	gitlab.com/contextualcode/platform_cc/api/router v0.0.1
	golang.org/x/crypto v0.0.0-20201124201722-c8d3bf9c5392
)
