package computer

import (
	"context"
	"fmt"
	"github.com/m4schini/computercraft-go/connection"
)

type Redstone interface {
	//SetOutput turns the redstone signal of a specific side on or off.
	SetOutput(ctx context.Context, side Side, on bool) error
	//SetAnalogOutput sets the redstone signal strength for a specific side.
	SetAnalogOutput(ctx context.Context, side Side, value int) error
	//Output gets the current redstone output of a specific side.
	Output(ctx context.Context, side Side) (bool, int, error)
	//Input gets the current redstone input of a specific side.
	Input(ctx context.Context, side Side) (int, error)
}

func SetOutput(conn connection.Connection, ctx context.Context, side Side, on bool) error {
	_, err := conn.Execute(ctx, fmt.Sprintf("")) //TODO
	return err
}

func SetAnalogOutput(conn connection.Connection, ctx context.Context, side Side, value int) error {
	_, err := conn.Execute(ctx, fmt.Sprintf("")) //TODO
	return err
}

func Output(conn connection.Connection, ctx context.Context, side Side) (bool, int, error) {
	_, err := conn.Execute(ctx, fmt.Sprintf("")) //TODO
	return false, 0, err
}

func Input(conn connection.Connection, ctx context.Context, side Side) (int, error) {
	_, err := conn.Execute(ctx, fmt.Sprintf("")) //TODO
	return 0, err
}
