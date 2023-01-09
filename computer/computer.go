package computer

import (
	"context"
	"github.com/m4schini/computercraft-go/computer/commands"
	"github.com/m4schini/computercraft-go/connection"
	"time"
)

type Computer interface {
	Shutdown(ctx context.Context) error
	Reboot(ctx context.Context) error
	Version(ctx context.Context) (string, error)

	ComputerId(ctx context.Context) (string, error)
	ComputerLabel(ctx context.Context) (string, error)
	SetComputerLabel(ctx context.Context, label string) error

	Uptime(ctx context.Context) (time.Duration, error)
	Time(ctx context.Context) (float64, error)

	IsTurtle(ctx context.Context) (bool, error)
	IsPocket(ctx context.Context) (bool, error)
}

type computer struct {
	conn connection.Connection
}

func NewComputer(conn connection.Connection) *computer {
	c := new(computer)
	c.conn = conn
	return c
}

func (c *computer) IsTurtle(ctx context.Context) (bool, error) {
	return connection.DoActionBool(ctx, c.conn, "turtle ~= nil")
}

func (c *computer) IsPocket(ctx context.Context) (bool, error) {
	return connection.DoActionBool(ctx, c.conn, "pocket ~= nil")
}

func (c *computer) Shutdown(ctx context.Context) error {
	conn := c.conn
	return commands.Shutdown(ctx, conn)
}

func (c *computer) Reboot(ctx context.Context) error {
	conn := c.conn
	return commands.Reboot(ctx, conn)
}

func (c *computer) Version(ctx context.Context) (string, error) {
	conn := c.conn
	return commands.Version(ctx, conn)
}

func (c *computer) ComputerId(ctx context.Context) (string, error) {
	conn := c.conn
	return commands.ComputerId(ctx, conn)
}

func (c *computer) ComputerLabel(ctx context.Context) (string, error) {
	conn := c.conn
	return commands.ComputerLabel(ctx, conn)
}

func (c *computer) SetComputerLabel(ctx context.Context, label string) error {
	conn := c.conn
	return commands.SetComputerLabel(ctx, conn, label)
}

func (c *computer) Uptime(ctx context.Context) (time.Duration, error) {
	conn := c.conn
	return commands.Uptime(ctx, conn)
}

func (c *computer) Time(ctx context.Context) (float64, error) {
	conn := c.conn
	return commands.Time(ctx, conn)
}
