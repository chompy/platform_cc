module gitlab.com/contextualcode/platform_cc/api/platformsh

go 1.16

replace gitlab.com/contextualcode/platform_cc/api/output => ../output

replace gitlab.com/contextualcode/platform_cc/api/def => ../def

replace gitlab.com/contextualcode/platform_cc/api/config => ../config

require (
	github.com/helloyi/go-sshclient v1.0.0
	github.com/melbahja/goph v1.2.1
	github.com/vaughan0/go-ini v0.0.0-20130923145212-a98ad7ee00ec
	github.com/ztrue/tracerr v0.3.0
	gitlab.com/contextualcode/platform_cc/api/config v0.0.1
	gitlab.com/contextualcode/platform_cc/api/output v0.0.1
	golang.org/x/oauth2 v0.0.0-20210413134643-5e61552d6c78
	golang.org/x/term v0.0.0-20201210144234-2321bbc49cbf
)
