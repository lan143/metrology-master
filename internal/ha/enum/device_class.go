package enum

type DeviceClass string

const (
	DeviceClassCurrent   DeviceClass = "current"
	DeviceClassEnergy    DeviceClass = "energy"
	DeviceClassFrequency DeviceClass = "frequency"
	DeviceClassPower     DeviceClass = "power"
	DeviceClassVoltage   DeviceClass = "voltage"
)
