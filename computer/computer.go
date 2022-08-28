package computer

import "time"

type Computer interface {
	Shutdown() error
	Reboot() error
	Version() (string, error)

	ComputerId() (string, error)
	ComputerLabel() (string, error)
	SetComputerLabel(label string) error

	Uptime() (time.Duration, error)
	Time() (float64, error)
}
