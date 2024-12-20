package registry

import (
	"encoding/json"
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
