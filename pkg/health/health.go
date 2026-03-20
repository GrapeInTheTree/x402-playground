package health

// Response represents a health check response.
type Response struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Network string `json:"network,omitempty"`
	Address string `json:"address,omitempty"`
}
