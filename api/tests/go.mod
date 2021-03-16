module gitlab.com/contextualcode/platform_cc/api/tests

go 1.15

replace gitlab.com/contextualcode/platform_cc/api/def => ../def

replace gitlab.com/contextualcode/platform_cc/api/project => ../project

replace gitlab.com/contextualcode/platform_cc/api/container => ../container

replace gitlab.com/contextualcode/platform_cc/api/router => ../router

replace gitlab.com/contextualcode/platform_cc/api/output => ../output

require (
	github.com/docker/docker v1.13.1
	github.com/ztrue/tracerr v0.3.0
	gitlab.com/contextualcode/platform_cc/api/container v0.0.1
	gitlab.com/contextualcode/platform_cc/api/def v0.0.1
	gitlab.com/contextualcode/platform_cc/api/output v0.0.1
	gitlab.com/contextualcode/platform_cc/api/project v0.0.1
	gitlab.com/contextualcode/platform_cc/api/router v0.0.1
	golang.org/x/lint v0.0.0-20201208152925-83fdc39ff7b5 // indirect
	golang.org/x/term v0.0.0-20210220032956-6a3ed077a48d // indirect
	golang.org/x/tools v0.0.0-20201217235154-5b06639e575e // indirect
)
