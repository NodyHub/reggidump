package registry

import (
	"encoding/json"
	"strings"
)

type Manifest struct {
	SchemaVersion int    `json:"schemaVersion"`
	Name          string `json:"name"`
	Tag           string `json:"tag"`
	Architecture  string `json:"architecture"`
	FsLayers      []struct {
		BlobSum string `json:"blobSum"`
	} `json:"fsLayers"`
	History []struct {
		V1Compatibility string `json:"v1Compatibility"`
		HistoryEntry    *HistoryEntry
	} `json:"history"`
	Signatures []struct {
		Header struct {
			Jwk struct {
				Crv string `json:"crv"`
				Kid string `json:"kid"`
				Kty string `json:"kty"`
				X   string `json:"x"`
				Y   string `json:"y"`
			} `json:"jwk"`
			Alg string `json:"alg"`
		} `json:"header"`
		Signature string `json:"signature"`
		Protected string `json:"protected"`
	} `json:"signatures"`
}

// ParseHistoryEntry parses the V1Compatibility field of a history entry
func (m *Manifest) ParseHistoryEntries() ([]*HistoryEntry, error) {
	// Parse the V1Compatibility field
	var historyEntries []*HistoryEntry
	for _, entry := range m.History {
		var entryStruct HistoryEntry
		if err := json.Unmarshal([]byte(entry.V1Compatibility), &entryStruct); err != nil {
			return nil, err
		}
		historyEntries = append(historyEntries, &entryStruct)
	}
	return historyEntries, nil
}

func GetFsLayerDotGraph(m *Manifest) string {
	if m == nil {
		return ""
	}
	if len(m.FsLayers) == 0 {
		return ""
	}
	if m.SchemaVersion == 2 {
		return getFsLayerDotGraphV2(m)
	}
	var sb strings.Builder
	sb.WriteString("\"" + m.FsLayers[0].BlobSum + "\"")
	for i := 1; i < len(m.FsLayers); i++ {
		sb.WriteString(" -> \"" + m.FsLayers[i].BlobSum + "\"")
	}
	sb.WriteString(";")
	return sb.String()
}

func getFsLayerDotGraphV2(m *Manifest) string {
	panic("not implemented")
}
