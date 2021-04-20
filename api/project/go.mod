module gitlab.com/contextualcode/platform_cc/api/project

go 1.15

replace gitlab.com/contextualcode/platform_cc/api/def => ../def

replace gitlab.com/contextualcode/platform_cc/api/container => ../container

replace gitlab.com/contextualcode/platform_cc/api/output => ../output

replace gitlab.com/contextualcode/platform_cc/api/platformsh => ../platformsh

replace gitlab.com/contextualcode/platform_cc/api/config => ../config

require (
	github.com/docker/docker v1.13.1
	github.com/martinlindhe/base36 v1.1.0
	github.com/ztrue/tracerr v0.3.0
	gitlab.com/contextualcode/platform_cc/api/config v0.0.1
	gitlab.com/contextualcode/platform_cc/api/container v0.0.1
	gitlab.com/contextualcode/platform_cc/api/def v0.0.1
	gitlab.com/contextualcode/platform_cc/api/output v0.0.1
	gitlab.com/contextualcode/platform_cc/api/platformsh v0.0.1
	golang.org/x/term v0.0.0-20210317153231-de623e64d2a6 // indirect
)
