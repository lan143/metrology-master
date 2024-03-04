package entity

import (
	"github.com/lan143/metrology-master/internal/ha/enum"
)

type Sensor struct {
	StateTopic        string           `json:"state_topic,omitempty"`
	ValueTemplate     string           `json:"value_template,omitempty"`
	UnitOfMeasurement string           `json:"unit_of_measurement,omitempty"`
	DeviceClass       enum.DeviceClass `json:"device_class,omitempty"`
	StateClass        enum.StateClass  `json:"state_class"`
	Base
}
