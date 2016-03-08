package main

import (
	"flag"
	"fmt"
	"github.com/kodokojo/kodokojo-haproxy-marathon/kodokojo"
	"net"
)

func main() {

	projectName := flag.String("projectName", "", "Project name to listen")
	port := flag.Int("httpPort", 8080, "port number to listen")
	marathonUrl := flag.String("marathonUrl", "http://localhost:8080", "Url to connect to Marathon API")
	marathonCallbackUrl := flag.String("marathonCallbackurl", "", "Marathon callback Url which will be registered on marathon")
	templatePath := flag.String("templatePath", "/Users/jpthiery/workspace/go/src/github.com/kodokojo/kodokojo-haproxy-marathon/kodokojo/default.tpl", "Path to the template file use to generate HA proxy configuration")

	flag.Parse()

	portStr := fmt.Sprintf(":%d", *port)

	if *marathonCallbackUrl == "" {

		ifaces, _ := net.Interfaces()

		var ip net.IP

		for _, i := range ifaces {
			addrs, _ := i.Addrs()
			for _, addr := range addrs {
				var tmp net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					tmp = v.IP
				case *net.IPAddr:
					tmp = v.IP
				}
				if !tmp.IsLoopback() {
					ip = tmp
				}
			}
		}
		*marathonCallbackUrl = "http://" + ip.String() + portStr + "/callback"
	}

	config := kodokojo.NewConfiguration(*projectName, *marathonUrl, *marathonCallbackUrl, *port, *templatePath)

	kodokojo.RegisterMarathon(config)

	handler := kodokojo.NewHttphandler(config)

	fmt.Println("Filtering project name		:", config.ProjectName())
	fmt.Println("Marathon url		:", config.MarathonUrl())
	fmt.Println("Marathon callback		:", config.MarathonCallbackUrl())

	generator := kodokojo.NewHaProxyConfigurationGenerator(*templatePath)

	locator := kodokojo.NewMarathonServiceLocator(config.MarathonUrl())
	var services kodokojo.Services = locator.LocateServices("acme", "ci")


	backends := make([]kodokojo.HaProxyBackEnd, 1)
	backends[0] = kodokojo.HaProxyBackEnd{BackEndHost:"1.2.3.4", BackEndPort:32749}

	haProxyEntries := make([]kodokojo.HaProxyEntry, 1)
	haProxyEntries[0] = kodokojo.HaProxyEntry{EntryName:"scm", Backends: backends}

	projects := make([]kodokojo.Project, 1)
	projects[0] = kodokojo.Project{ProjectName:"acme", SSHIp:"192.168.99.100", SSHPort:22022, HaProxyHTTPEntries: services.HaProxyHTTPEntries, HaProxySSHEntries:haProxyEntries}

	context := kodokojo.HaProxyContext{Projects: projects, SyslogEntryPoint:"/dev/log"}

	output := generator.GenerateConfiguration(context)
	fmt.Println("Configuration generated :",output)

	handler.Start()

}
