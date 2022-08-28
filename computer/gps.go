package computer

import "time"

type GPS interface {
	Locate() (int, int, int, error)
	LocateWithTimeout(timeout time.Duration) (int, int, int, error)
}
