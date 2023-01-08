package computer

import (
	"context"
	"errors"
	"fmt"
	"github.com/m4schini/computercraft-go/computer/commands"
	"github.com/m4schini/computercraft-go/connection"
	"go.uber.org/multierr"
	"io"
	"time"
)

type Turtle interface {
	Forward(ctx context.Context) (bool, error)
	Back(ctx context.Context) (bool, error)
	Up(ctx context.Context) (bool, error)
	Down(ctx context.Context) (bool, error)
	TurnLeft(ctx context.Context) (bool, error)
	TurnRight(ctx context.Context) (bool, error)

	Dig(ctx context.Context) (bool, error)
	DigDown(ctx context.Context) (bool, error)
	DigUp(ctx context.Context) (bool, error)

	Place(ctx context.Context) (bool, error)
	PlaceUp(ctx context.Context) (bool, error)
	PlaceDown(ctx context.Context) (bool, error)

	Drop(ctx context.Context, count int) (bool, error)
	DropUp(ctx context.Context, count int) (bool, error)
	DropDown(ctx context.Context, count int) (bool, error)

	//Select changes the currently selected slot.
	Select(ctx context.Context, slot int) (bool, error)
	//SelectedSlot gets the currently selected slot.
	SelectedSlot(ctx context.Context) (int, error)
	//ItemCount gets the number of items in the given slot.
	ItemCount(ctx context.Context, slot int) (int, error)
	//ItemSpace gets the remaining number of items which may be stored in this stack.
	ItemSpace(ctx context.Context, slot int) (int, error)
	//ItemDetail gets detailed information about the items in the given slot.
	ItemDetail(ctx context.Context, slot int, detailed bool) (map[string]interface{}, error)
	//CompareTo compares the item in the currently selected slot to the item in another slot.
	CompareTo(ctx context.Context, slot int) (bool, error)
	//TransferTo moves an item from the selected slot to another one.
	TransferTo(ctx context.Context, slot, count int) (bool, error)

	//Detect checks if there is a solid block in front of the turtle. In this case,
	//solid refers to any non-air or liquid block.
	Detect(ctx context.Context) bool
	DetectUp(ctx context.Context) bool
	DetectDown(ctx context.Context) bool

	//Compare checks if the block in front of the turtle is equal to the item in the currently selected slot.
	Compare(ctx context.Context) (bool, error)
	CompareUp(ctx context.Context) (bool, error)
	CompareDown(ctx context.Context) (bool, error)

	//Attack attacks the entity in front of the turtle.
	Attack(ctx context.Context) (bool, error)
	AttackUp(ctx context.Context) (bool, error)
	AttackDown(ctx context.Context) (bool, error)

	//Suck sucks an item from the inventory in front of the turtle, or from an item floating in the world.
	Suck(ctx context.Context, count int) (bool, error)
	SuckUp(ctx context.Context, count int) (bool, error)
	SuckDown(ctx context.Context, count int) (bool, error)

	//FuelLevel gets the maximum amount of fuel this turtle currently holds.
	FuelLevel(ctx context.Context) (int, error)
	//Refuel refuels this turtle
	Refuel(ctx context.Context, count int) (bool, error)
	//FuelLimit gets the maximum amount of fuel this turtle can hold.
	FuelLimit(ctx context.Context) (int, error)

	//Inspect gets information about the block in front of the turtle.
	Inspect(ctx context.Context) (bool, Block, error)
	InspectUp(ctx context.Context) (bool, Block, error)
	InspectDown(ctx context.Context) (bool, Block, error)

	//Craft crafts a recipe based on the turtle's inventory.
	Craft(ctx context.Context, limit int) (bool, error)

	io.Closer
	Computer
	Settings
	GPS
}

type turtle struct {
	conn connection.Connection
}

func NewTurtle(conn connection.Connection) *turtle {
	t := new(turtle)
	t.conn = conn
	return t
}

func (t *turtle) IsTurtle(ctx context.Context) (bool, error) {
	return connection.DoActionBool(ctx, t.conn, "turtle ~= nil")
}

func (t *turtle) IsPocket(ctx context.Context) (bool, error) {
	return connection.DoActionBool(ctx, t.conn, "pocket ~= nil")
}

func (t *turtle) Close() error {
	ctx, stop := context.WithTimeout(context.Background(), 1*time.Second)
	defer stop()

	return multierr.Combine(
		t.Shutdown(ctx),
		t.conn.Close(),
	)
}

func (t *turtle) Shutdown(ctx context.Context) error {
	conn := t.conn
	return commands.Shutdown(ctx, conn)
}

func (t *turtle) Reboot(ctx context.Context) error {
	conn := t.conn
	return commands.Reboot(ctx, conn)
}

func (t *turtle) Version(ctx context.Context) (string, error) {
	conn := t.conn
	return commands.Version(ctx, conn)
}

func (t *turtle) ComputerId(ctx context.Context) (string, error) {
	conn := t.conn
	return commands.ComputerId(ctx, conn)
}

func (t *turtle) ComputerLabel(ctx context.Context) (string, error) {
	conn := t.conn
	return commands.ComputerLabel(ctx, conn)
}

func (t *turtle) SetComputerLabel(ctx context.Context, label string) error {
	conn := t.conn
	return commands.SetComputerLabel(ctx, conn, label)
}

func (t *turtle) Uptime(ctx context.Context) (time.Duration, error) {
	conn := t.conn
	return commands.Uptime(ctx, conn)
}

func (t *turtle) Time(ctx context.Context) (float64, error) {
	conn := t.conn
	return commands.Time(ctx, conn)
}

func (t *turtle) Forward(ctx context.Context) (bool, error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, "turtle.forward()")
}

func (t *turtle) Back(ctx context.Context) (bool, error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, "turtle.back()")
}

func (t *turtle) Up(ctx context.Context) (bool, error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, "turtle.up()")
}

func (t *turtle) Down(ctx context.Context) (bool, error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, "turtle.down()")
}

func (t *turtle) TurnLeft(ctx context.Context) (bool, error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, "turtle.turnLeft()")
}

func (t *turtle) TurnRight(ctx context.Context) (bool, error) {
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

func (t *turtle) Detect(ctx context.Context) bool {
	return t._detect(ctx, "turtle.detect()")
}

func (t *turtle) DetectUp(ctx context.Context) bool {
	return t._detect(ctx, "turtle.detectUp()")

}

func (t *turtle) DetectDown(ctx context.Context) bool {
	return t._detect(ctx, "turtle.detectDown()")

}

func (t *turtle) Dig(ctx context.Context) (bool, error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, "turtle.dig()")
}

func (t *turtle) DigDown(ctx context.Context) (bool, error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, "turtle.digDown()")
}

func (t *turtle) DigUp(ctx context.Context) (bool, error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, "turtle.digUp()")
}

func (t *turtle) Place(ctx context.Context) (bool, error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, "turtle.place()")
}

func (t *turtle) PlaceUp(ctx context.Context) (bool, error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, "turtle.placeUp()")

}

func (t *turtle) PlaceDown(ctx context.Context) (bool, error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, "turtle.placeDown()")

}

func (t *turtle) Drop(ctx context.Context, count int) (bool, error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, fmt.Sprintf("turtle.drop(%v)", count))

}

func (t *turtle) DropUp(ctx context.Context, count int) (bool, error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, fmt.Sprintf("turtle.dropUp(%v)", count))
}

func (t *turtle) DropDown(ctx context.Context, count int) (bool, error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, fmt.Sprintf("turtle.dropDown(%v)", count))
}

func (t *turtle) Select(ctx context.Context, slot int) (bool, error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, fmt.Sprintf("turtle.select(%v)", slot))
}

func (t *turtle) SelectedSlot(ctx context.Context) (int, error) {
	conn := t.conn
	return connection.DoActionInt(ctx, conn, "turtle.getSelectedSlot()")
}

func (t *turtle) ItemCount(ctx context.Context, slot int) (int, error) {
	conn := t.conn
	return connection.DoActionInt(ctx, conn, fmt.Sprintf("turtle.getItemCount(%v)", slot))
}

func (t *turtle) ItemSpace(ctx context.Context, slot int) (int, error) {
	conn := t.conn
	return connection.DoActionInt(ctx, conn, fmt.Sprintf("turtle.getItemSpace(%v)", slot))
}

func (t *turtle) ItemDetail(ctx context.Context, slot int, detailed bool) (map[string]interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (t *turtle) CompareTo(ctx context.Context, slot int) (bool, error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, fmt.Sprintf("turtle.compareTo(%v)", slot))
}

func (t *turtle) TransferTo(ctx context.Context, slot, count int) (bool, error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, fmt.Sprintf("turtle.transferTo(%v,%v)", slot, count))
}

func (t *turtle) Compare(ctx context.Context) (bool, error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, "turtle.compare()")
}

func (t *turtle) CompareUp(ctx context.Context) (bool, error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, "turtle.compareUp()")

}

func (t *turtle) CompareDown(ctx context.Context) (bool, error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, "turtle.compareDown()")

}

func (t *turtle) Attack(ctx context.Context) (bool, error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, "turtle.attack()")

}

func (t *turtle) AttackUp(ctx context.Context) (bool, error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, "turtle.attackUp()")

}

func (t *turtle) AttackDown(ctx context.Context) (bool, error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, "turtle.attackDown()")

}

func (t *turtle) Suck(ctx context.Context, count int) (bool, error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, fmt.Sprintf("turtle.suck(%v)", count))
}

func (t *turtle) SuckUp(ctx context.Context, count int) (bool, error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, fmt.Sprintf("turtle.suckUp(%v)", count))
}

func (t *turtle) SuckDown(ctx context.Context, count int) (bool, error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, fmt.Sprintf("turtle.suckDown(%v)", count))
}

func (t *turtle) FuelLevel(ctx context.Context) (int, error) {
	conn := t.conn
	return connection.DoActionInt(ctx, conn, "turtle.getFuelLevel()")
}

func (t *turtle) Refuel(ctx context.Context, count int) (bool, error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, fmt.Sprintf("turtle.refuel(%v)", count))
}

func (t *turtle) FuelLimit(ctx context.Context) (int, error) {
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

func (t *turtle) Inspect(ctx context.Context) (bool, Block, error) {
	return t._doInspect(ctx, "turtle.inspect()")
}

func (t *turtle) InspectUp(ctx context.Context) (bool, Block, error) {
	return t._doInspect(ctx, "turtle.inspectUp()")

}

func (t *turtle) InspectDown(ctx context.Context) (bool, Block, error) {
	return t._doInspect(ctx, "turtle.inspectDown()")

}

func (t *turtle) Craft(ctx context.Context, limit int) (bool, error) {
	conn := t.conn
	return connection.DoActionBool(ctx, conn, fmt.Sprintf("turtle.craft(%v)", limit))
}

func (t *turtle) _doLocate(ctx context.Context, timeout time.Duration, debug bool) (int, int, int, error) {
	conn := t.conn
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

func (t *turtle) Locate(ctx context.Context) (int, int, int, error) {
	return t._doLocate(ctx, 2*time.Second, false)
}

func (t *turtle) LocateWithTimeout(ctx context.Context, timeout time.Duration) (int, int, int, error) {
	return t._doLocate(ctx, timeout, false)
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
