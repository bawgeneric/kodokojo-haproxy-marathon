package kodokojo


type Configuration struct {
	projectName string
	marathonUrl string
	marathonCallbackUrl string
	port int
	templatePath string
}

func NewConfiguration(projectName string,marathonUrl string, marathonCallbackUrl string,  port int, templatePath string) Configuration {
	return Configuration{projectName, marathonUrl,marathonCallbackUrl,  port, templatePath}
}

func (c *Configuration) Port() int {
	return c.port
}

func (c *Configuration) ProjectName() string {
	return c.projectName;
}

func (c *Configuration) MarathonCallbackUrl() string {
	return c.marathonCallbackUrl;
}

func (c *Configuration) MarathonUrl() string {
	return c.marathonUrl;
}

func (c *Configuration) TemplatePath() string {
	return c.templatePath;
}