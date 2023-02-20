package computer

import (
	"context"
	"errors"
	"fmt"
	"github.com/m4schini/computercraft-go/computer/gps"
	"github.com/m4schini/computercraft-go/connection"
	"time"
)

type turtle struct {
	conn connection.Connection
}

func NewTurtle(conn connection.Connection) *turtle {
	t := new(turtle)
	t.conn = conn
	return t
}

func (t *turtle) IsTurtle(ctx context.Context) (isTurtle bool, err error) {
	return connection.DoActionBool(ctx, t.conn, "turtle ~= nil")
}

func (t *turtle) IsPocket(ctx context.Context) (isPocket bool, err error) {
	return connection.DoActionBool(ctx, t.conn, "pocket ~= nil")
}

func (t *turtle) Shutdown(ctx context.Context) error {
	return Shutdown(ctx, t.conn)
}

func (t *turtle) Reboot(ctx context.Context) error {
	return Reboot(ctx, t.conn)
}

func (t *turtle) Version(ctx context.Context) (version string, err error) {
	return Version(ctx, t.conn)
}

func (t *turtle) ComputerId(ctx context.Context) (id string, err error) {
	return ComputerId(ctx, t.conn)
}

func (t *turtle) ComputerLabel(ctx context.Context) (label string, err error) {
	return ComputerLabel(ctx, t.conn)
}

func (t *turtle) SetComputerLabel(ctx context.Context, label string) error {
	return SetComputerLabel(ctx, t.conn, label)
}

func (t *turtle) Uptime(ctx context.Context) (uptime time.Duration, err error) {
	return Uptime(ctx, t.conn)
}

func (t *turtle) Time(ctx context.Context) (time float64, err error) {
	return Time(ctx, t.conn)
}

func (t *turtle) Forward(ctx context.Context) (success bool, err error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, "turtle.forward()")
}

func (t *turtle) Back(ctx context.Context) (success bool, err error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, "turtle.back()")
}

func (t *turtle) Up(ctx context.Context) (success bool, err error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, "turtle.up()")
}

func (t *turtle) Down(ctx context.Context) (success bool, err error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, "turtle.down()")
}

func (t *turtle) TurnLeft(ctx context.Context) (success bool, err error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, "turtle.turnLeft()")
}

func (t *turtle) TurnRight(ctx context.Context) (success bool, err error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, "turtle.turnRight()")
}

func (t *turtle) _detect(ctx context.Context, command string) bool {
	conn := t.conn
	res, err := conn.Execute(ctx, command)
	if err != nil {
		return false
	}

	detect, ok := res[0].(bool)
	return ok && detect
}

func (t *turtle) Detect(ctx context.Context) (detected bool) {
	return t._detect(ctx, "turtle.detect()")
}

func (t *turtle) DetectUp(ctx context.Context) (detected bool) {
	return t._detect(ctx, "turtle.detectUp()")

}

func (t *turtle) DetectDown(ctx context.Context) (detected bool) {
	return t._detect(ctx, "turtle.detectDown()")

}

func (t *turtle) Dig(ctx context.Context) (success bool, err error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, "turtle.dig()")
}

func (t *turtle) DigDown(ctx context.Context) (success bool, err error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, "turtle.digDown()")
}

func (t *turtle) DigUp(ctx context.Context) (success bool, err error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, "turtle.digUp()")
}

func (t *turtle) Place(ctx context.Context) (placed bool, err error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, "turtle.place()")
}

func (t *turtle) PlaceUp(ctx context.Context) (placed bool, err error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, "turtle.placeUp()")

}

func (t *turtle) PlaceDown(ctx context.Context) (placed bool, err error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, "turtle.placeDown()")

}

func (t *turtle) Drop(ctx context.Context, count int) (dropped bool, err error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, fmt.Sprintf("turtle.drop(%v)", count))

}

func (t *turtle) DropUp(ctx context.Context, count int) (dropped bool, err error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, fmt.Sprintf("turtle.dropUp(%v)", count))
}

func (t *turtle) DropDown(ctx context.Context, count int) (dropped bool, err error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, fmt.Sprintf("turtle.dropDown(%v)", count))
}

func (t *turtle) Select(ctx context.Context, slot int) (selected bool, err error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, fmt.Sprintf("turtle.select(%v)", slot))
}

func (t *turtle) SelectedSlot(ctx context.Context) (slot int, err error) {
	conn := t.conn
	return connection.DoActionInt(ctx, conn, "turtle.getSelectedSlot()")
}

func (t *turtle) ItemCount(ctx context.Context, slot int) (count int, err error) {
	conn := t.conn
	return connection.DoActionInt(ctx, conn, fmt.Sprintf("turtle.getItemCount(%v)", slot))
}

func (t *turtle) ItemSpace(ctx context.Context, slot int) (space int, err error) {
	conn := t.conn
	return connection.DoActionInt(ctx, conn, fmt.Sprintf("turtle.getItemSpace(%v)", slot))
}

func (t *turtle) ItemDetail(ctx context.Context, slot int, detailed bool) (map[string]interface{}, error) {
	conn := t.conn
	response, err := conn.Execute(ctx, fmt.Sprintf(`turtle.getItemDetail(%v, %v)`, slot, detailed))
	if err != nil {
		return nil, err
	}

	if len(response) < 1 {
		return nil, errors.New("unexpected data length")
	}

	slotdata, ok := response[0].(map[string]interface{})
	if !ok {
		return nil, connection.UnexpectedDatatypeErr
	}

	return slotdata, nil
}

func (t *turtle) CompareTo(ctx context.Context, slot int) (same bool, err error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, fmt.Sprintf("turtle.compareTo(%v)", slot))
}

func (t *turtle) TransferTo(ctx context.Context, slot, count int) (transferred bool, err error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, fmt.Sprintf("turtle.transferTo(%v,%v)", slot, count))
}

func (t *turtle) Compare(ctx context.Context) (same bool, err error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, "turtle.compare()")
}

func (t *turtle) CompareUp(ctx context.Context) (same bool, err error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, "turtle.compareUp()")

}

func (t *turtle) CompareDown(ctx context.Context) (same bool, err error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, "turtle.compareDown()")

}

func (t *turtle) Attack(ctx context.Context) (success bool, err error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, "turtle.attack()")

}

func (t *turtle) AttackUp(ctx context.Context) (success bool, err error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, "turtle.attackUp()")

}

func (t *turtle) AttackDown(ctx context.Context) (success bool, err error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, "turtle.attackDown()")

}

func (t *turtle) Suck(ctx context.Context, count int) (success bool, err error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, fmt.Sprintf("turtle.suck(%v)", count))
}

func (t *turtle) SuckUp(ctx context.Context, count int) (success bool, err error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, fmt.Sprintf("turtle.suckUp(%v)", count))
}

func (t *turtle) SuckDown(ctx context.Context, count int) (success bool, err error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, fmt.Sprintf("turtle.suckDown(%v)", count))
}

func (t *turtle) FuelLevel(ctx context.Context) (fuelLevel int, err error) {
	conn := t.conn
	return connection.DoActionInt(ctx, conn, "turtle.getFuelLevel()")
}

func (t *turtle) Refuel(ctx context.Context, count int) (success bool, err error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, fmt.Sprintf("turtle.refuel(%v)", count))
}

func (t *turtle) FuelLimit(ctx context.Context) (fuelLimit int, err error) {
	conn := t.conn
	return connection.DoActionInt(ctx, conn, "turtle.getFuelLimit()")
}

func (t *turtle) _doInspect(ctx context.Context, command string) (bool, Block, error) {
	conn := t.conn
	res, err := conn.Execute(ctx, command)
	if err != nil {
		return false, nil, connection.RpcError(err)
	}

	if len(res) < 2 {
		return false, nil, connection.RpcError(errors.New("not enough parameter"))
	}

	errMsg, isError := res[1].(string)
	if isError {
		return false, nil, connection.RpcError(errors.New(errMsg))
	}

	detectedBlock, ok := res[0].(bool)
	if !detectedBlock {
		return false, nil, nil
	}

	data, ok := res[1].(map[string]interface{})
	if !ok {
		return false, nil, connection.UnexpectedDatatypeErr
	}

	return detectedBlock, data, nil
}

func (t *turtle) Inspect(ctx context.Context) (detected bool, block Block, err error) {
	return t._doInspect(ctx, "turtle.inspect()")
}

func (t *turtle) InspectUp(ctx context.Context) (detected bool, block Block, err error) {
	return t._doInspect(ctx, "turtle.inspectUp()")

}

func (t *turtle) InspectDown(ctx context.Context) (detected bool, block Block, err error) {
	return t._doInspect(ctx, "turtle.inspectDown()")

}

func (t *turtle) Craft(ctx context.Context, limit int) (success bool, err error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, fmt.Sprintf("turtle.craft(%v)", limit))
}

func (t *turtle) Locate(ctx context.Context) (x int, y int, z int, err error) {
	return gps.Locate(ctx, t.conn)
}

func (t *turtle) LocateWithTimeout(ctx context.Context, timeout time.Duration) (x int, y int, z int, err error) {
	return gps.LocateWithTimeout(ctx, t.conn, timeout)
}

func (t *turtle) Define(ctx context.Context, name string, option ...SettingsOption) error {
	conn := t.conn
	_, err := conn.Execute(ctx, fmt.Sprintf("settings.define(\"%s\")", name))
	return err
}

func (t *turtle) Undefine(ctx context.Context, name string) error {
	conn := t.conn
	_, err := conn.Execute(ctx, fmt.Sprintf("settings.undefine(\"%s\")", name))
	return err
}

func (t *turtle) Set(ctx context.Context, name, value string) error {
	conn := t.conn
	_, err := conn.Execute(ctx, fmt.Sprintf("settings.set(\"%s\", \"%s\")", name, value))
	return err
}

func (t *turtle) Unset(ctx context.Context, name string) error {
	conn := t.conn
	_, err := conn.Execute(ctx, fmt.Sprintf("settings.unset(\"%s\")", name))
	return err
}

func (t *turtle) Get(ctx context.Context, name string) (string, error) {
	conn := t.conn
	res, err := conn.Execute(ctx, fmt.Sprintf("settings.set(\"%s\")", name))
	if err != nil {
		return "", err
	}
	str, ok := res[0].(string)
	if !ok {
		return "", errors.New("not a string")
	}
	return str, nil
}

func (t *turtle) Clear(ctx context.Context) error {
	conn := t.conn
	_, err := conn.Execute(ctx, fmt.Sprintf("settings.clear()"))
	return err
}

func (t *turtle) Names(ctx context.Context) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (t *turtle) Load(ctx context.Context, path string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (t *turtle) Save(ctx context.Context, path string) (bool, error) {
	//TODO implement me
	panic("implement me")
}
