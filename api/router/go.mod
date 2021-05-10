module gitlab.com/contextualcode/platform_cc/router

go 1.15

replace gitlab.com/contextualcode/platform_cc/api/project => ../project

replace gitlab.com/contextualcode/platform_cc/api/def => ../def

replace gitlab.com/contextualcode/platform_cc/api/output => ../output

replace gitlab.com/contextualcode/platform_cc/api/container => ../container

replace gitlab.com/contextualcode/platform_cc/api/platformsh => ../platformsh

replace gitlab.com/contextualcode/platform_cc/api/config => ../config

require (
	github.com/pkg/errors v0.9.1
	gitlab.com/contextualcode/platform_cc/api/config v0.0.1
	gitlab.com/contextualcode/platform_cc/api/container v0.0.1
	gitlab.com/contextualcode/platform_cc/api/def v0.0.1
	gitlab.com/contextualcode/platform_cc/api/output v0.0.1
	gitlab.com/contextualcode/platform_cc/api/project v0.0.1
)
