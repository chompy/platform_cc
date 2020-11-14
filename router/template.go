package router

const nginxTemplate = `
{{ range . }}
server {
    resolver 127.0.0.11;
    server_name {{ .host }};
    listen 80;
    client_max_body_size 200M;
    {{ range .routes }}
	{{ if eq .type "upstream" }}
	location {{ .path }} {
		{{ range $k, $v := .route.Redirects.Paths }}
		location {{ $k }} {
			return {{ $v.Code }} {{ $v.To }};
		}
		{{ end }}
		location ~* {
			proxy_pass http://{{ .upstream }};
			proxy_set_header X-Client-IP $server_addr;
			proxy_set_header X-Forwarded-Host $host;
			proxy_set_header X-Forwarded-Port $server_port;
			proxy_set_header X-Forwarded-Proto $scheme;
			proxy_set_header X-Forwarded-Server $host;
			proxy_set_header Host $host;
			proxy_set_header X-Forwarded-For $remote_addr;
		}
	}
    {{ else if eq .type "redirect" }}
    location {{ .path }} {
        return 301 {{ .to }};
    }
    {{ end }}
    {{ end }}
}
{{ end }}
`
