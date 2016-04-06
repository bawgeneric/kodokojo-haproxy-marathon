package commons

import "time"

type HaProxyContext struct {
	Projects         []Project
	SyslogEntryPoint string
}

func (h *HaProxyContext) AddProject(project Project) {
	h.Projects = append(h.Projects, project)
}

type Project struct {
	ProjectName        string
	LastChangedAd      time.Time
	Version            time.Time
	LastConfigChangeAt time.Time
	LastScalingAt      time.Time
	SSHPort            int
	HaProxyHTTPEntries []HaProxyEntry
	HaProxySSHEntries  []HaProxyEntry
}


func (p *Project) IsReady() bool {
	return len(p.HaProxyHTTPEntries) > 0 || len(p.HaProxySSHEntries) > 0
}

func (p *Project) IsHTTPReady() bool {
	return p.haveBackend(p.HaProxyHTTPEntries)
}

func (p *Project) IsSSHReady() bool {
	return p.haveBackend(p.HaProxySSHEntries)
}

func (p *Project) haveBackend(entries []HaProxyEntry) (res bool) {
	for _, entry := range entries {
		res = len(entry.Backends) > 0
		if !res {
			break
		}
	}
	return
}

type HaProxyEntry struct {
	EntryName      string
	EntryComponent string
	Backends       []HaProxyBackEnd
}

type HaProxyBackEnd struct {
	BackEndHost string
	BackEndPort int
}

type Service struct {
	ProjectName        string
	Version            time.Time
	LastConfigChangeAt time.Time
	LastScalingAt      time.Time
	HaProxyHTTPEntries []HaProxyEntry
	HaProxySSHEntries  []HaProxyEntry
}

type MarathonEvent struct {
	EventType   string    `json:"eventType"`
	AppId       string    `json:"appId"`
	Timestamp   time.Time `json:"timestamp"`
	Alive       bool      `json:"alive"`
	CallbackUrl string    `json:"callbackUrl"`
}
