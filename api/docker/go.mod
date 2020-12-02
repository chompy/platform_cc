module gitlab.com/contextualcode/platform_cc/api/docker

go 1.15

replace gitlab.com/contextualcode/platform_cc/api/def => ../def

replace gitlab.com/contextualcode/platform_cc/api/output => ../output

require (
	github.com/Microsoft/go-winio v0.4.15 // indirect
	github.com/creack/pty v1.1.11
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.13.1
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-units v0.4.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/ztrue/tracerr v0.3.0
	golang.org/x/crypto v0.0.0-20201124201722-c8d3bf9c5392
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b // indirect
	gitlab.com/contextualcode/platform_cc/api/output v0.0.1
)
