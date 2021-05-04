module gitlab.com/contextualcode/platform_cc/api/config

go 1.15

replace gitlab.com/contextualcode/platform_cc/api/output => ../output

replace gitlab.com/contextualcode/platform_cc/api/def => ../def

require (
	github.com/ztrue/tracerr v0.3.0
	gitlab.com/contextualcode/platform_cc/api/def v0.0.1
	gitlab.com/contextualcode/platform_cc/api/output v0.0.1
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b
	golang.org/x/sys v0.0.0-20201218084310-7d0127a74742 // indirect
	golang.org/x/term v0.0.0-20201210144234-2321bbc49cbf // indirect
)
