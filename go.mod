module gitlab.com/contextualcode/platform_cc

go 1.15

replace gitlab.com/contextualcode/platform_cc/api => ./api
replace gitlab.com/contextualcode/platform_cc/cmd => ./cmd
replace gitlab.com/contextualcode/platform_cc/def => ./def

require (
	github.com/fatih/color v1.10.0 // indirect
	github.com/goccy/go-yaml v1.8.3 // indirect
	github.com/gopherjs/vecty v0.5.0 // indirect
	github.com/gorhill/cronexpr v0.0.0-20180427100037-88b0669f7d75
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/nsf/termbox-go v0.0.0-20200418040025-38ba6e5628f1
	github.com/spf13/cobra v1.1.1
	gitlab.com/contextualcode/platform_cc/api v0.0.1
	gitlab.com/contextualcode/platform_cc/cmd v0.0.1
	golang.org/x/sys v0.0.0-20201110211018-35f3e6cf4a65 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/yaml.v2 v2.3.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
)
