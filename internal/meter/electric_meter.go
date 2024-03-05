package meter

import "context"

type ElectricMeter interface {
	GetPowerConsumption(ctx context.Context) (float64, error)
	GetFrequency(ctx context.Context) (float64, error)
	GetVoltage(ctx context.Context) (float64, error)
	GetCurrent(ctx context.Context) (float64, error)
	GetActivePower(ctx context.Context) (float64, error)
	GetReactivePower(ctx context.Context) (float64, error)
	GetFullPower(ctx context.Context) (float64, error)

	Meter
}
