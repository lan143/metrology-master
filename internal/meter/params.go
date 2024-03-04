package meter

type Flags uint64

const (
	FlagHasPowerConsumption = 0x01
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
