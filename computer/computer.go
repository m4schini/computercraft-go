package computer

import (
	"github.com/m4schini/computercraft-go/computer/commands"
	"github.com/m4schini/computercraft-go/connection"
	"time"
)

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

type computer struct {
	client connection.Client
}

func NewComputer(client connection.Client) *computer {
	c := new(computer)
	c.client = client
	return c
}

func (c *computer) Shutdown() error {
	conn := c.client.Connection()
	return commands.Shutdown(conn.Context(), conn)
}

func (c *computer) Reboot() error {
	conn := c.client.Connection()
	return commands.Reboot(conn.Context(), conn)
}

func (c *computer) Version() (string, error) {
	conn := c.client.Connection()
	return commands.Version(conn.Context(), conn)
}

func (c *computer) ComputerId() (string, error) {
	conn := c.client.Connection()
	return commands.ComputerId(conn.Context(), conn)
}

func (c *computer) ComputerLabel() (string, error) {
	conn := c.client.Connection()
	return commands.ComputerLabel(conn.Context(), conn)
}

func (c *computer) SetComputerLabel(label string) error {
	conn := c.client.Connection()
	return commands.SetComputerLabel(conn.Context(), conn, label)
}

func (c *computer) Uptime() (time.Duration, error) {
	conn := c.client.Connection()
	return commands.Uptime(conn.Context(), conn)
}

func (c *computer) Time() (float64, error) {
	conn := c.client.Connection()
	return commands.Time(conn.Context(), conn)
}

func (c *computer) Close() error {
	return c.client.Connection().Close()
}