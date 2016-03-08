global
maxconn 4096
log {{ .SyslogEntryPoint }}   local0
log {{ .SyslogEntryPoint }}   local1 info

defaults

option  dontlognull

timeout connect 5000ms
timeout client 50000ms
timeout server 50000ms

frontend http-in
  log     global
  mode    http
  bind *:80
  reqadd X-Forwarded-Proto:\ http

frontend https-in
  log     global
  mode    http
  bind *:443 ssl{{range .Projects}}{{range .HaProxyHTTPEntries}} crt/usr/local/etc/haproxy/ssl/{{.EntryName}}-server.pem {{end}}{{end}}
  reqadd X-Forwarded-Proto:\ https

  option httplog
  option dontlognull
  option forwardfor
  option http-server-close

{{range .Projects}}# BEGIN entries project {{.ProjectName}}{{$projectName := .ProjectName}}{{range .HaProxyHTTPEntries}}
  acl host_{{$projectName}}_{{.EntryName}} hdr_beg(host) -i {{.EntryName}}.{{$projectName}}{{end}}{{range .HaProxyHTTPEntries}}
  use_backend {{.EntryName}}-{{$projectName}}-cluster-http if host_{{$projectName}}_{{.EntryName}}{{end}}
# END entries project {{.ProjectName}}{{end}}

  stats enable
  stats auth admin:admin
  stats uri /stats

{{range .Projects}}{{$projectName := .ProjectName}}{{$sshPort := .SSHPort}}frontend ssh-{{.ProjectName}}-in
  bind    *:{{.SSHPort}}
  default_backend {{$projectName}}-cluster-ssh

{{range .HaProxyHTTPEntries}}backend {{.EntryName}}-{{$projectName}}-cluster-http
  mode    http
  redirect scheme https if !{ ssl_fc }
  balance leastconn{{$entryName := .EntryName}}
  {{range $index,$backend := .Backends}}server {{$entryName}}{{$projectName}}{{$index}} {{$backend.BackEndHost}}:{{$backend.BackEndPort}} check{{end}}
{{end}}
{{range .HaProxySSHEntries}}backend {{$projectName}}-cluster-ssh{{$entryName := .EntryName}}
  {{range $index,$backend := .Backends}}server {{$entryName}}{{$projectName}}{{$index}} {{$backend.BackEndHost}}:{{$backend.BackEndPort}} check port {{$backend.BackEndPort}}{{end}}
{{end}}{{end}}

