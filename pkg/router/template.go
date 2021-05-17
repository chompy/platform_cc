/*
This file is part of Platform.CC.

Platform.CC is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Platform.CC is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Platform.CC.  If not, see <https://www.gnu.org/licenses/>.
*/

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
    server {
        server_name default;
        listen 80 default;
        listen 443 ssl;
        ssl_certificate /var/ssl/localhost/cert.pem;
        ssl_certificate_key /var/ssl/localhost/key.pem;
        root /www;
        location / {
            index index.html;
        }
    }
    include /routes/*.conf;
}
`

const nginxServerTemplate = `
{{ range . }}
server {
    resolver 127.0.0.11;
    server_name {{ .host }};
    listen 80;
    listen 443 ssl;
    ssl_certificate /var/ssl/localhost/cert.pem;
    ssl_certificate_key /var/ssl/localhost/key.pem;
    client_max_body_size 200M;
    {{ range .routes }}
	{{ if eq .type "upstream" }}
	location "{{ .path }}" {
		{{ range .redirects }}
		location ~ "{{ .path }}" {
			return {{ .code }} {{ .to }};
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
    location "{{ .path }}" {
        return 301 {{ .to }};
    }
    {{ end }}
    {{ end }}
}
{{ end }}
`

const routeListHTML = `
<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>Platform.CC Routes</title>
    <style type="text/css">
        html, body { font-family: sans-serif; }
        table { width: 90%; }
        th, td { padding: 5px; }
        th { text-align: left; background-color: teal; }
    </style>
</head>
<body>
    <div id="container">
        <h1>Routes</h1>
        <div id="inject">Loading...</div>
    </div>
    <script type="text/javascript">
        function fetchFile(name, callback) {
            fetch(name)
                .then(response => response.text()) 
                .then(textString => {
                    callback(textString);
                }
            );            
        }
        function createTableHead(values) {
            let tr = document.createElement("thead");
            for (let i in values) {
                let td = document.createElement("th")
                td.innerText = values[i];
                tr.appendChild(td);
            }
            return tr;
        }
        function createTableRow(values) {
            let tr = document.createElement("tr");
            for (let i in values) {
                let td = document.createElement("td")
                td.innerHTML = values[i];
                tr.appendChild(td);
            }
            return tr;
        }
        function addProject(id, data) {
            let e = document.createElement("div");
            let header = document.createElement("h2");
            header.innerText = id;
            e.appendChild(header);
            let t = document.createElement("table");
            t.appendChild(
                createTableHead(["Host", "Upstream"])
            );
            for (let i in data) {
                let upstream = "";
                for (let j in data[i].routes) {
                    if (data[i].routes[j].type == "upstream") {
                        upstream = data[i].routes[j].upstream;
                    }
                }
                t.appendChild(
                    createTableRow(["<a target='_blank' href='//"+data[i].host+"'>"+data[i].host+"</a>", upstream])
                );
            }
            e.appendChild(t);
            document.getElementById("inject").appendChild(e);
        }
        fetchFile("projects.txt", function(data) {
            data = data.split("\n");
            let processed = [];
            document.getElementById("inject").innerHTML = "No routes found.";
            for (let i in data) {
                let fileName = data[i];
                if (!fileName || processed.indexOf(fileName) > -1) {
                    continue;
                }
                if (processed.length == 0) {
                    document.getElementById("inject").innerHTML = "";
                }
                fetchFile(fileName + ".json", function(data) {
                    addProject(fileName, JSON.parse(data));
                });
                processed.push(fileName);
            }
        });
    </script>
</body>
</html>
`
