package meter

type Flags uint64

const (
	FlagHasPowerConsumption = 0x01
	FlagHasFrequency        = 0x02
	FlagHasVoltage          = 0x04
	FlagHasCurrent          = 0x08
	FlagHasActivePower      = 0x10
	FlagHasReactivePower    = 0x20
	FlagHasFullPower        = 0x040
)

type Params struct {
	UID          string
	StateTopic   string
	Manufacturer string
	Model        string
	Name         string
	HWVersion    string
	SWVersion    string
	Flags        Flags
}

func (f Flags) HasPowerConsumption() bool {
	return f&FlagHasPowerConsumption > 0
}

func (f Flags) HasFrequency() bool {
	return f&FlagHasFrequency > 0
}

func (f Flags) HasVoltage() bool {
	return f&FlagHasVoltage > 0
}

func (f Flags) HasCurrent() bool {
	return f&FlagHasCurrent > 0
}

func (f Flags) HasActivePower() bool {
	return f&FlagHasActivePower > 0
}

func (f Flags) HasReactivePower() bool {
	return f&FlagHasReactivePower > 0
}

func (f Flags) HasFullPower() bool {
	return f&FlagHasFullPower > 0
}
