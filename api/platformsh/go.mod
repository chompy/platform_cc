module gitlab.com/contextualcode/platform_cc/api/platformsh

go 1.16

replace gitlab.com/contextualcode/platform_cc/api/output => ../output

replace gitlab.com/contextualcode/platform_cc/api/def => ../def

require (
	github.com/martinlindhe/base36 v1.1.0 // indirect
	github.com/melbahja/goph v1.2.1 // indirect
	github.com/vaughan0/go-ini v0.0.0-20130923145212-a98ad7ee00ec
	github.com/ztrue/tracerr v0.3.0
	gitlab.com/contextualcode/platform_cc/api/def v0.0.1
	gitlab.com/contextualcode/platform_cc/api/output v0.0.1
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2 // indirect
)