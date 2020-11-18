package router

const nginxBaseConf = `
user nginx;
worker_processes 1;
error_log  /var/log/nginx/error.log warn;
pid        /var/run/nginx.pid;
events {
    worker_connections  1024;
}
http {
    server_names_hash_bucket_size 128;
    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 65;
    types_hash_max_size 2048;
    server_tokens off;
    include /etc/nginx/mime.types;
    types {
        application/vnd.apple.pkpass pkpass;
    }
    default_type application/octet-stream;
    absolute_redirect off;
    proxy_request_buffering off;
    fastcgi_request_buffering off;
    client_max_body_size 0;
    client_body_buffer_size 128k;
    proxy_buffering on;
    fastcgi_buffering on;
    proxy_buffer_size 32k;
    fastcgi_buffer_size 32k;
    proxy_buffers 128 4k;
    fastcgi_buffers 128 4k;
    proxy_busy_buffers_size 32k;
    fastcgi_busy_buffers_size 32k;
    proxy_max_temp_file_size 0;
    fastcgi_max_temp_file_size 0;
    proxy_connect_timeout     30s;
    fastcgi_connect_timeout   30s;
    proxy_read_timeout        86400s;
    fastcgi_read_timeout      86400s;
    proxy_send_timeout        86400s;
    fastcgi_send_timeout      86400s;
    include /routes/*.conf;
}
`

const nginxServerTemplate = `
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
		location ~ {{ $k }} {
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
