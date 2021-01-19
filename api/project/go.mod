module gitlab.com/contextualcode/platform_cc/api/project

go 1.15

replace gitlab.com/contextualcode/platform_cc/api/def => ../def

replace gitlab.com/contextualcode/platform_cc/api/container => ../container

replace gitlab.com/contextualcode/platform_cc/api/output => ../output

require (
	github.com/martinlindhe/base36 v1.1.0
	github.com/ztrue/tracerr v0.3.0
	gitlab.com/contextualcode/platform_cc/api/def v0.0.1
	gitlab.com/contextualcode/platform_cc/api/container v0.0.1
	gitlab.com/contextualcode/platform_cc/api/output v0.0.1
)
