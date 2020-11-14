module gitlab.com/contextualcode/platform_cc/cmd

go 1.15

replace gitlab.com/contextualcode/platform_cc/api => ../api

replace gitlab.com/contextualcode/platform_cc/def => ../def

require (
	github.com/spf13/cobra v1.1.1
	gitlab.com/contextualcode/platform_cc/api v0.0.1
	gitlab.com/contextualcode/platform_cc/def v0.0.1
)
