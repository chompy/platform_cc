module gitlab.com/contextualcode/platform_cc

go 1.15

replace gitlab.com/contextualcode/platform_cc/api/project => ./api/project

replace gitlab.com/contextualcode/platform_cc/api/container => ./api/container

replace gitlab.com/contextualcode/platform_cc/api/output => ./api/output

replace gitlab.com/contextualcode/platform_cc/api/def => ./api/def

replace gitlab.com/contextualcode/platform_cc/api/router => ./api/router

replace gitlab.com/contextualcode/platform_cc/cli => ./cli

replace gitlab.com/contextualcode/platform_cc/api/tests => ./api/tests

require (
	gitlab.com/contextualcode/platform_cc/api/output v0.0.1
	gitlab.com/contextualcode/platform_cc/cli v0.0.1

)
