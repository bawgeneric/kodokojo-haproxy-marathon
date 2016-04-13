package commons

type Configuration struct {
	projectNameMatch    string
	haProxyCfgPath      string
	marathonUrl         string
	marathonCallbackUrl string
	port                int
	templatePath        string
}

func (c *Configuration) Port() int {
	return c.port
}

func (c *Configuration) getProjectNameMatch() string {
	//Now, dead code. May be used to dedicate an HA proxy to a given project/entity
	return c.projectNameMatch
}

func (c *Configuration) HaProxyCfgPath() string {
	return c.haProxyCfgPath
}

func (c *Configuration) MarathonCallbackUrl() string {
	return c.marathonCallbackUrl
}

func (c *Configuration) MarathonUrl() string {
	return c.marathonUrl
}

func (c *Configuration) TemplatePath() string {
	return c.templatePath
}
