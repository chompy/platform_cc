module gitlab.com/contextualcode/platform_cc/cmd

go 1.15

replace gitlab.com/contextualcode/platform_cc/api/project => ../api/project

replace gitlab.com/contextualcode/platform_cc/api/docker => ../api/docker

replace gitlab.com/contextualcode/platform_cc/api/def => ../api/def

replace gitlab.com/contextualcode/platform_cc/api/router => ../api/router

require (
	github.com/spf13/cobra v1.1.1
	github.com/ztrue/tracerr v0.3.0
	gitlab.com/contextualcode/platform_cc/api/def v0.0.1
	gitlab.com/contextualcode/platform_cc/api/docker v0.0.1
	gitlab.com/contextualcode/platform_cc/api/project v0.0.1
	gitlab.com/contextualcode/platform_cc/api/router v0.0.1
)
