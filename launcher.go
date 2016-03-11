package main

import (
	"flag"
	"fmt"
	"github.com/kodokojo/kodokojo-haproxy-marathon/commons"
	"github.com/kodokojo/kodokojo-haproxy-marathon/haproxy"
	"github.com/kodokojo/kodokojo-haproxy-marathon/marathon"
	"github.com/kodokojo/kodokojo-haproxy-marathon/utils"
	"log"
	"net"
)

func main() {

	port := flag.Int("httpPort", 8080, "port number to listen")
	marathonUrl := flag.String("marathonUrl", "http://localhost:8080", "Url to connect to Marathon API")
	haProxyCfgPath := flag.String("haProxyCfgPath", "/Users/jpthiery/workspace/go/src/github.com/kodokojo/kodokojo-haproxy-marathon//haproxy.cfg", "HaProxy.cfg configuration Path")
	marathonCallbackUrl := flag.String("marathonCallbackUrl", "", "Marathon callback Url which will be registered on marathon")
	templatePath := flag.String("templatePath", "", "Path to the template file use to generate HA proxy configuration")
	projectName := flag.String("projectName", "", "Project name to listen - Not used in this version.")

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

	marathonEventChannel := make(chan commons.MarathonEvent, 5)

	config := commons.NewConfiguration(*projectName, *haProxyCfgPath, *marathonUrl, *marathonCallbackUrl, *port, *templatePath)
	locator := marathon.NewMarathonServiceLocator(config.MarathonUrl())
	sslStore := utils.NewSslStore(*marathonUrl)
	generator := haproxy.NewHaProxyConfigurator(*templatePath, sslStore)
	applicationState := haproxy.NewApplicationsState(config, locator, generator, commons.HaProxyContext{})
	applicationState.Start(marathonEventChannel)

	marathon.RegisterMarathon(config)

	handler := marathon.NewHttphandler(config, marathonEventChannel)

	log.Println("Marathon url		:", config.MarathonUrl())
	log.Println("Marathon callback		:", config.MarathonCallbackUrl())

	services, _ := locator.LocateAllService()
	for _, service := range services {
		applicationState.UpdateIfConfigurationChanged(service)
	}

	handler.Start()

}
