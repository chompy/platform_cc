module gitlab.com/contextualcode/platform_cc/router

go 1.15

replace gitlab.com/contextualcode/platform_cc/api => ../api

replace gitlab.com/contextualcode/platform_cc/def => ../def

require (
	gitlab.com/contextualcode/platform_cc/api v0.0.1
	gitlab.com/contextualcode/platform_cc/def v0.0.1
)
