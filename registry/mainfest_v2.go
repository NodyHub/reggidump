package registry

import "time"

type ManifestV2 struct {
	SchemaVersion int    `json:"schemaVersion"`
	MediaType     string `json:"mediaType"`
	Config        struct {
		MediaType string `json:"mediaType"`
		Size      int    `json:"size"`
		Digest    string `json:"digest"`
	} `json:"config"`
	Layers []struct {
		MediaType string `json:"mediaType"`
		Size      int    `json:"size"`
		Digest    string `json:"digest"`
	} `json:"layers"`
}

type Config struct {
	ExposedPorts map[string]struct{} `json:"ExposedPorts"`
	Env          []string            `json:"Env"`
	Entrypoint   []string            `json:"Entrypoint"`
	Cmd          []string            `json:"Cmd"`
	WorkingDir   string              `json:"WorkingDir"`
	ArgsEscaped  bool                `json:"ArgsEscaped"`
}

type History struct {
	Created    time.Time `json:"created"`
	CreatedBy  string    `json:"created_by"`
	Comment    string    `json:"comment,omitempty"`
	EmptyLayer bool      `json:"empty_layer,omitempty"`
}

type Rootfs struct {
	Type    string   `json:"type"`
	DiffIDs []string `json:"diff_ids"`
}

type DockerImage struct {
	Architecture string    `json:"architecture"`
	Config       Config    `json:"config"`
	Created      time.Time `json:"created"`
	History      []History `json:"history"`
	OS           string    `json:"os"`
	Rootfs       Rootfs    `json:"rootfs"`
}
