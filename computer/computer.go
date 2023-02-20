package computer

import (
	"context"
	"errors"
	"fmt"
	"github.com/m4schini/computercraft-go/connection"
	"time"
)

func Shutdown(ctx context.Context, conn connection.Connection) error {
	_, err := conn.Execute(ctx, "os.shutdown()")
	return err
}

func Reboot(ctx context.Context, conn connection.Connection) error {
	_, err := conn.Execute(ctx, "os.reboot()")
	return err
}

func Version(ctx context.Context, conn connection.Connection) (string, error) {
	res, err := conn.Execute(ctx, "os.version()")
	if err != nil {
		return "", err
	}

	version, ok := res[0].(string)
	if ok {
		return version, nil
	} else {
		return "", connection.UnexpectedDatatypeErr
	}
}

func ComputerId(ctx context.Context, conn connection.Connection) (string, error) {
	res, err := conn.Execute(ctx, "os.getComputerID()")
	if err != nil {
		return "", connection.RpcError(err)
	}

	if len(res) != 1 {
		return "", errors.New("something went wrong")
	}

	label, ok := res[0].(float64)
	if ok {
		return fmt.Sprintf("%v", label), nil
	}

	return "", connection.UnexpectedDatatypeErr
}

func ComputerLabel(ctx context.Context, conn connection.Connection) (string, error) {
	res, err := conn.Execute(ctx, "os.getComputerLabel()")
	if err != nil {
		return "", connection.RpcError(err)
	}

	if len(res) != 1 {
		return "", errors.New("something went wrong")
	}

	label, ok := res[0].(string)
	if ok {
		return fmt.Sprintf("%v", label), nil
	}

	return "", connection.UnexpectedDatatypeErr
}

func SetComputerLabel(ctx context.Context, conn connection.Connection, label string) error {
	var err error
	if label == "" {
		_, err = conn.Execute(ctx, "os.setComputerLabel()")
	} else {
		_, err = conn.Execute(ctx, fmt.Sprintf("os.setComputerLabel(\"%v\")", label))
	}

	if err != nil {
		return connection.RpcError(err)
	} else {
		return nil
	}
}

func Uptime(ctx context.Context, conn connection.Connection) (time.Duration, error) {
	res, err := conn.Execute(ctx, "os.clock()")
	if err != nil {
		return 0, connection.RpcError(err)
	}

	uptime, ok := res[0].(float64)
	if !ok {
		return 0, connection.UnexpectedDatatypeErr
	}

	return time.Duration(uptime) * time.Second, nil
}

func Time(ctx context.Context, conn connection.Connection) (float64, error) {
	res, err := conn.Execute(ctx, "os.time()")
	if err != nil {
		return 0, connection.RpcError(err)
	}

	tme, ok := res[0].(float64)
	if !ok {
		return 0, connection.RpcError(err)
	}

	return tme, nil
}

type computer struct {
	conn connection.Connection
}

func NewComputer(conn connection.Connection) *computer {
	c := new(computer)
	c.conn = conn
	return c
}

func (c *computer) Shutdown(ctx context.Context) error {
	return Shutdown(ctx, c.conn)
}

func (c *computer) Reboot(ctx context.Context) error {
	return Reboot(ctx, c.conn)
}

func (c *computer) Version(ctx context.Context) (version string, err error) {
	return Version(ctx, c.conn)
}

func (c *computer) ComputerId(ctx context.Context) (id string, err error) {
	return ComputerId(ctx, c.conn)
}

func (c *computer) ComputerLabel(ctx context.Context) (label string, err error) {
	return ComputerLabel(ctx, c.conn)
}

func (c *computer) SetComputerLabel(ctx context.Context, label string) error {
	return SetComputerLabel(ctx, c.conn, label)
}

func (c *computer) Uptime(ctx context.Context) (uptime time.Duration, err error) {
	return Uptime(ctx, c.conn)
}

func (c *computer) Time(ctx context.Context) (t float64, err error) {
	return Time(ctx, c.conn)
}
