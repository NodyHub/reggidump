package config

type Common struct {
	Auth      string   `json:"auth,omitempty"`
	Filter    string   `json:"filter,omitempty"`
	Header    []string `json:"header,omitempty"`
	LogLevel  string   `json:"log_level,omitempty"`
	LogFile   string   `json:"log_file,omitempty"`
	Rertry    string   `json:"retry,omitempty"`
	Timeout   int      `json:"timeout,omitempty"`
	UserAgent string   `json:"user_agent,omitempty"`
}

type List struct {
	Common `json:",inline"`
	Output string `json:"output,omitempty"`
}

type Dump struct {
	Common       `json:",inline"`
	Destination  string `json:"destination,omitempty"`
	ManifestOnly bool   `json:"manifest_only,omitempty"`
	Parallel     int    `json:"parallel,omitempty"`
}
