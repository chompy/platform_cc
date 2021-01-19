module gitlab.com/contextualcode/platform_cc/router

go 1.15

replace gitlab.com/contextualcode/platform_cc/api/project => ../project

replace gitlab.com/contextualcode/platform_cc/api/def => ../def

replace gitlab.com/contextualcode/platform_cc/api/output => ../output

replace gitlab.com/contextualcode/platform_cc/api/container => ../container

require (
	github.com/ztrue/tracerr v0.3.0
	gitlab.com/contextualcode/platform_cc/api/project v0.0.1
	gitlab.com/contextualcode/platform_cc/api/def v0.0.1
	gitlab.com/contextualcode/platform_cc/api/output v0.0.1
	gitlab.com/contextualcode/platform_cc/api/container v0.0.1
)
