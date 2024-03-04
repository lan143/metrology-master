package meter

import "context"

type ElectricMeter interface {
	GetPowerConsumption(ctx context.Context) (float64, error)

	Meter
}
