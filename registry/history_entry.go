package registry

type HistoryEntry struct {
	Architecture string `json:"architecture,omitempty"`
	Config       struct {
		Hostname     string      `json:"Hostname,omitempty"`
		Domainname   string      `json:"Domainname,omitempty"`
		User         string      `json:"User,omitempty"`
		AttachStdin  bool        `json:"AttachStdin,omitempty"`
		AttachStdout bool        `json:"AttachStdout,omitempty"`
		AttachStderr bool        `json:"AttachStderr,omitempty"`
		Digest       string      `json:"Digest,omitempty"`
		Tty          bool        `json:"Tty,omitempty"`
		OpenStdin    bool        `json:"OpenStdin,omitempty"`
		StdinOnce    bool        `json:"StdinOnce,omitempty"`
		Env          []string    `json:"Env,omitempty"`
		Cmd          []string    `json:"Cmd,omitempty"`
		Image        string      `json:"Image,omitempty"`
		Volumes      interface{} `json:"Volumes,omitempty"`
		WorkingDir   string      `json:"WorkingDir,omitempty"`
		Entrypoint   []string    `json:"Entrypoint,omitempty"`
		OnBuild      interface{} `json:"OnBuild,omitempty"`
		Labels       struct {
			OrgOpencontainersImageRefName string `json:"org.opencontainers.image.ref.name,omitempty"`
			OrgOpencontainersImageVersion string `json:"org.opencontainers.image.version,omitempty"`
		} `json:"Labels,omitempty"`
	} `json:"config,omitempty"`
	Container       string `json:"container,omitempty"`
	ContainerConfig struct {
		Hostname     string      `json:"Hostname,omitempty"`
		Domainname   string      `json:"Domainname,omitempty"`
		User         string      `json:"User,omitempty"`
		AttachStdin  bool        `json:"AttachStdin,omitempty"`
		AttachStdout bool        `json:"AttachStdout,omitempty"`
		AttachStderr bool        `json:"AttachStderr,omitempty"`
		Tty          bool        `json:"Tty,omitempty"`
		OpenStdin    bool        `json:"OpenStdin,omitempty"`
		StdinOnce    bool        `json:"StdinOnce,omitempty"`
		Env          []string    `json:"Env,omitempty"`
		Cmd          []string    `json:"Cmd,omitempty"`
		Image        string      `json:"Image,omitempty"`
		Volumes      interface{} `json:"Volumes,omitempty"`
		WorkingDir   string      `json:"WorkingDir,omitempty"`
		Entrypoint   interface{} `json:"Entrypoint,omitempty"`
		OnBuild      interface{} `json:"OnBuild,omitempty"`
		Labels       struct {
			OrgOpencontainersImageRefName string `json:"org.opencontainers.image.ref.name,omitempty"`
			OrgOpencontainersImageVersion string `json:"org.opencontainers.image.version,omitempty"`
		} `json:"Labels,omitempty"`
	} `json:"container_config,omitempty"`
	Created       string `json:"created,omitempty"`
	DockerVersion string `json:"docker_version,omitempty"`
	ID            string `json:"id,omitempty"`
	OS            string `json:"os,omitempty"`
	Parent        string `json:"parent,omitempty"`
	Throwaway     bool   `json:"throwaway,omitempty"`
	Variant       string `json:"variant,omitempty"`
}

func (h *HistoryEntry) Clone() *HistoryEntry {
	clone := *h
	clone.Config = h.Config
	return &clone
}
