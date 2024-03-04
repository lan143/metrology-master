package entity

type Device struct {
	ConfigurationUrl string   `json:"configuration_url,omitempty"`
	Connections      []string `json:"connections,omitempty"`
	HWVersion        string   `json:"hw_version,omitempty"`
	Identifiers      []string `json:"identifiers,omitempty"`
	Manufacturer     string   `json:"manufacturer,omitempty"`
	Model            string   `json:"model,omitempty"`
	Name             string   `json:"name,omitempty"`
	SuggestedArea    string   `json:"suggested_area,omitempty"`
	SWVersion        string   `json:"sw_version,omitempty"`
	ViaDevice        string   `json:"viaDevice,omitempty"`
}
