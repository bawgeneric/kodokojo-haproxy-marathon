package kodokojo

import (
	"text/template"
	"bytes"
	"os"
)

type HaProxyContext struct {
	Projects []Project
	SyslogEntryPoint string
}

type Project struct {
	ProjectName string
	SSHIp string
	SSHPort int
	HaProxyHTTPEntries []HaProxyEntry
	HaProxySSHEntries []HaProxyEntry
}

type HaProxyEntry struct {
	EntryName string
	Backends []HaProxyBackEnd
}

type HaProxyBackEnd struct {
	BackEndHost string
	BackEndPort int
}

type haProxyConfigurationGenerator struct {
	templatePath string
}

func NewHaProxyConfigurationGenerator(templatePath string) haProxyConfigurationGenerator {
	return haProxyConfigurationGenerator{templatePath}
}

func (g *haProxyConfigurationGenerator) GenerateConfiguration(context HaProxyContext) string {
	tmpl, err := template.New("default.tpl").ParseFiles(g.templatePath)
	if err != nil {
		panic(err)
	}
	var writer bytes.Buffer
	errExe := tmpl.Execute(os.Stdout, context)
	if errExe != nil {
		panic(errExe)
	}
	return writer.String()
}