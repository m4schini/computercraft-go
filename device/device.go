package device

import "github.com/m4schini/computercraft-go/connection"

type Capability string

const (
	CapabilityGPS    Capability = "gps"
	CapabilityPocket Capability = "pocket"
)

type Device interface {
}

type device struct {
	conn connection.Connection
}
