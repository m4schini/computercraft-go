package gps

import (
	"context"
	"errors"
	"fmt"
	"github.com/m4schini/computercraft-go/connection"
	"time"
)

type Vector struct {
	x, y, z int
}

const (
	HeadingNorth Heading = iota
	HeadingEast
	HeadingSouth
	HeadingWest
	HeadingUp
	HeadingDown
)

type Heading uint8

func (h Heading) Vector() Vector {
	switch h {
	case HeadingNorth:
		return Vector{x: 0, y: 0, z: -1}
	case HeadingEast:
		return Vector{x: +1, y: 0, z: 0}
	case HeadingSouth:
		return Vector{x: 0, y: 0, z: +1}
	case HeadingWest:
		return Vector{x: -1, y: 0, z: 0}
	case HeadingUp:
		return Vector{x: 0, y: +1, z: 0}
	case HeadingDown:
		return Vector{x: 0, y: -1, z: 0}
	default:
		panic("unknown heading")
	}
}

const (
	DefaultTimeout = 2 * time.Second
)

type GPS interface {
	Locate(ctx context.Context) (int, int, int, error)
	LocateWithTimeout(ctx context.Context, timeout time.Duration) (int, int, int, error)
}

func _doLocate(ctx context.Context, conn connection.Connection, timeout time.Duration, debug bool) (int, int, int, error) {
	res, err := conn.Execute(
		ctx,
		fmt.Sprintf("gps.locate(%v, %v)", int(timeout.Seconds()), debug),
	)
	if err != nil {
		return 0, 0, 0, connection.RpcError(err)
	}

	if len(res) == 0 {
		return 0, 0, 0, errors.New("position could not be established")
	}

	if len(res) < 3 {
		return 0, 0, 0, errors.New("position could not be established")
	}

	x, ok := res[0].(float64)
	if !ok {
		return 0, 0, 0, connection.UnexpectedDatatypeErr
	}

	y, ok := res[1].(float64)
	if !ok {
		return 0, 0, 0, connection.UnexpectedDatatypeErr
	}

	z, ok := res[2].(float64)
	if !ok {
		return 0, 0, 0, connection.UnexpectedDatatypeErr
	}

	return int(x), int(y), int(z), nil
}

func Locate(ctx context.Context, conn connection.Connection) (int, int, int, error) {
	return _doLocate(ctx, conn, 2*time.Second, false)
}

func LocateWithTimeout(ctx context.Context, conn connection.Connection, timeout time.Duration) (int, int, int, error) {
	return _doLocate(ctx, conn, timeout, false)
}

type ModemGPS struct {
	Conn connection.Connection
}

func (m *ModemGPS) Locate(ctx context.Context) (int, int, int, error) {
	return _doLocate(ctx, m.Conn, DefaultTimeout, false)
}

func (m *ModemGPS) LocateWithTimeout(ctx context.Context, timeout time.Duration) (int, int, int, error) {
	return _doLocate(ctx, m.Conn, timeout, false)
}

type relativeGPS struct {
	x, y, z int
	heading Heading
}

func (r *relativeGPS) Locate(ctx context.Context) (int, int, int, error) {
	return r.x, r.y, r.z, nil
}

func (r *relativeGPS) LocateWithTimeout(ctx context.Context, timeout time.Duration) (int, int, int, error) {
	return r.x, r.y, r.z, nil
}
