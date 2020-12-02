module gitlab.com/contextualcode/platform_cc/api/tests

go 1.15

replace gitlab.com/contextualcode/platform_cc/api/def => ../def
replace gitlab.com/contextualcode/platform_cc/api/project => ../project
replace gitlab.com/contextualcode/platform_cc/api/docker => ../docker
replace gitlab.com/contextualcode/platform_cc/api/router => ../router
replace gitlab.com/contextualcode/platform_cc/api/output => ../output

require (
	github.com/ztrue/tracerr v0.3.0
	gitlab.com/contextualcode/platform_cc/api/def v0.0.1
	gitlab.com/contextualcode/platform_cc/api/project v0.0.1
	gitlab.com/contextualcode/platform_cc/api/docker v0.0.1
	gitlab.com/contextualcode/platform_cc/api/router v0.0.1
	gitlab.com/contextualcode/platform_cc/api/output v0.0.1
)
