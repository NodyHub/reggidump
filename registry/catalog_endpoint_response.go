package registry

type CatalogEndpointResponse struct {
	Repositories []string `json:"repositories,omitempty"`
}
