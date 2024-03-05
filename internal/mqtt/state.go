package mqtt

type State struct {
	PowerConsumption float64 `json:"powerConsumption"`
	Frequency        float64 `json:"frequency"`
	Voltage          float64 `json:"voltage"`
	Current          float64 `json:"current"`
	ActivePower      float64 `json:"activePower"`
	ReactivePower    float64 `json:"reactivePower"`
	FullPower        float64 `json:"fullPower"`
}
