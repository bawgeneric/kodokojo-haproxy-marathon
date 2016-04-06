package haproxy

import (
	"bytes"
	"github.com/kodokojo/kodokojo-haproxy-marathon/commons"
	"github.com/kodokojo/kodokojo-haproxy-marathon/utils"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"text/template"
)

const sslPath string = "/usr/local/etc/haproxy/ssl/"

type haProxyConfigurator struct {
	templatePath string
	sslStore     utils.SslStore
	cache        map[string][]byte
}

func NewHaProxyConfigurator(templatePath string, sslStore utils.SslStore) haProxyConfigurator {
	return haProxyConfigurator{templatePath, sslStore, make(map[string][]byte)}
}

func (g *haProxyConfigurator) GenerateConfiguration(context commons.HaProxyContext) string {
	var tpl template.Template
	if len(g.templatePath) > 0 {
		tmplRead := template.Must(template.ParseFiles(g.templatePath))
		tpl = *tmplRead
	} else {
		tmplRead := template.Must(template.New("default.tpl").Parse(defaultHaProxyTemplate))
		tpl = *tmplRead
	}
	var writer bytes.Buffer
	errExe := tpl.Execute(&writer, context)
	if errExe != nil {
		panic(errExe)
	}
	return writer.String()
}

func (h *haProxyConfigurator) ReloadHaProxyWithConfiguration(haConfiguration string, configuration commons.Configuration, haProxyContext commons.HaProxyContext) {
	log.Println("Reloading configuration")
	log.Println(&haProxyContext)
	if len(haConfiguration) > 0 {
		err := ioutil.WriteFile(configuration.HaProxyCfgPath(), []byte(haConfiguration), 0644)
		if err != nil {
			log.Fatal("Unable to write file", configuration.HaProxyCfgPath(), err)
		}
	}

	for _, project := range haProxyContext.Projects {
		for _, entry := range project.HaProxyHTTPEntries {
			key := h.generateKey(project.ProjectName, entry.EntryName)
			if _, exist := h.cache[key]; !exist {
				sskKeyPath := h.generateCertFilePath(project.ProjectName, entry.EntryName)
				if stat, err := os.Stat(sskKeyPath); os.IsNotExist(err) || stat.Size() <= 0 {
					sslKeyContent := h.sslStore.GetPemFileFromSslStore(project.ProjectName, entry.EntryName)
					err := ioutil.WriteFile(sskKeyPath, sslKeyContent, 0644)
					h.cache[key] = sslKeyContent
					log.Println("Retrive a SSL certificate for key", key, "write in file", sskKeyPath)
					log.Println(sskKeyPath, "content", string(sslKeyContent))
					if err != nil {
						log.Println("Unable to write SSL file", sskKeyPath, err)
					}
				}
			} else {
				log.Println("Found an existing SSL certificate for key", key)
			}
		}
	}
	pidData, _ := ioutil.ReadFile("/tmp/haproxy.pid")
	pid := string(pidData)

	reload := exec.Command("/usr/local/sbin/haproxy", "-D", "-f", configuration.HaProxyCfgPath(), "-p", "/tmp/haproxy.pid", "-sf", pid)

	err := reload.Run()
	if err != nil {
		output, _ := reload.CombinedOutput()
		log.Println("Error while trying to reload HA proxy wiht command'", "/usr/local/sbin/haproxy -D -f "+configuration.HaProxyCfgPath()+" -p /tmp/haproxy.pid -sf "+pid, "' :", err)
		log.Println(string(output))
		//		log.Println(errOutput.Error())
	}
}

func (h *haProxyConfigurator) generateKey(projectName string, entityType string) string {
	return projectName + "-" + entityType
}

func (h *haProxyConfigurator) generateCertFilePath(projectName string, entityType string) string {
	return sslPath + projectName + "-" + entityType + "-server.pem"
}

const defaultHaProxyTemplate string = `global
	maxconn 4096
	log 192.168.99.100 local0
	log 192.168.99.100 local1 info

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
{{range .Projects}}{{if .IsHTTPReady}}# BEGIN entries project {{.ProjectName}}{{$projectName := .ProjectName}}{{range .HaProxyHTTPEntries}}
	acl host_{{$projectName}}_{{.EntryName}} hdr_beg(host) -i {{.EntryName}}.{{$projectName}}{{end}}{{range .HaProxyHTTPEntries}}
	use_backend {{.EntryName}}-{{$projectName}}-cluster-http if host_{{$projectName}}_{{.EntryName}}{{end}}
# END entries project {{.ProjectName}}{{end}}{{end}}

frontend https-in
	log     global
	mode    http
	bind *:443 ssl{{range .Projects}}{{$projectName := .ProjectName}}{{if .IsReady}}{{range .HaProxyHTTPEntries}} crt /usr/local/etc/haproxy/ssl/{{$projectName}}-{{.EntryName}}-server.pem {{end}}{{end}}{{end}}
	reqadd X-Forwarded-Proto:\ https

	option httplog
	option dontlognull
	option forwardfor
	option http-server-close

{{range .Projects}}{{if .IsHTTPReady}}# BEGIN entries project {{.ProjectName}}{{$projectName := .ProjectName}}{{range .HaProxyHTTPEntries}}
	acl host_{{$projectName}}_{{.EntryName}} hdr_beg(host) -i {{.EntryName}}.{{$projectName}}{{end}}{{range .HaProxyHTTPEntries}}
	use_backend {{.EntryName}}-{{$projectName}}-cluster-http if host_{{$projectName}}_{{.EntryName}}{{end}}
# END entries project {{.ProjectName}}{{end}}{{end}}

	stats enable
	stats auth admin:admin
	stats uri /stats
{{range .Projects}}{{$projectName := .ProjectName}}{{if .IsReady}}
{{if .IsSSHReady}}{{$sshPort := .SSHPort}}frontend ssh-{{.ProjectName}}-in
	bind    *:{{.SSHPort}}
	default_backend {{$projectName}}-cluster-ssh

{{range .HaProxySSHEntries}}backend {{$projectName}}-cluster-ssh{{$entryName := .EntryName}}
	{{range $index,$backend := .Backends}}server {{$entryName}}{{$projectName}}{{$index}} {{$backend.BackEndHost}}:{{$backend.BackEndPort}} check port {{$backend.BackEndPort}}{{end}}
{{end}}{{end}}

{{if .IsHTTPReady}}{{range .HaProxyHTTPEntries}}backend {{.EntryName}}-{{$projectName}}-cluster-http
	mode    http
	redirect scheme https if !{ ssl_fc }
	balance leastconn{{$entryName := .EntryName}}
	{{range $index,$backend := .Backends}}server {{$entryName}}{{$projectName}}{{$index}} {{$backend.BackEndHost}}:{{$backend.BackEndPort}} check
{{end}}{{end}}{{end}}
{{end}}{{end}}`
