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
	github.com/Microsoft/go-winio v0.4.16 // indirect
	github.com/olekukonko/tablewriter v0.0.4 // indirect
	gitlab.com/contextualcode/platform_cc/api/output v0.0.1
	gitlab.com/contextualcode/platform_cc/cli v0.0.1
	gitlab.com/contextualcode/platform_cc/cmd v0.0.0-20210303212337-fc25da2c0f60 // indirect
	golang.org/x/term v0.0.0-20210220032956-6a3ed077a48d // indirect

)
