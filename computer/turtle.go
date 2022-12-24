package computer

import (
	"errors"
	"fmt"
	"github.com/m4schini/computercraft-go/computer/commands"
	"github.com/m4schini/computercraft-go/connection"
	"io"
	"time"
)

type Turtle interface {
	Forward() (bool, error)
	Back() (bool, error)
	Up() (bool, error)
	Down() (bool, error)
	TurnLeft() (bool, error)
	TurnRight() (bool, error)

	Dig() (bool, error)
	DigDown() (bool, error)
	DigUp() (bool, error)

	Place() (bool, error)
	PlaceUp() (bool, error)
	PlaceDown() (bool, error)

	Drop(count int) (bool, error)
	DropUp(count int) (bool, error)
	DropDown(count int) (bool, error)

	//Select changes the currently selected slot.
	Select(slot int) (bool, error)
	//SelectedSlot gets the currently selected slot.
	SelectedSlot() (int, error)
	//ItemCount gets the number of items in the given slot.
	ItemCount(slot int) (int, error)
	//ItemSpace gets the remaining number of items which may be stored in this stack.
	ItemSpace(slot int) (int, error)
	//ItemDetail gets detailed information about the items in the given slot.
	ItemDetail(slot int, detailed bool) (map[string]interface{}, error)
	//CompareTo compares the item in the currently selected slot to the item in another slot.
	CompareTo(slot int) (bool, error)
	//TransferTo moves an item from the selected slot to another one.
	TransferTo(slot, count int) (bool, error)

	//Detect checks if there is a solid block in front of the turtle. In this case,
	//solid refers to any non-air or liquid block.
	Detect() bool
	DetectUp() bool
	DetectDown() bool

	//Compare checks if the block in front of the turtle is equal to the item in the currently selected slot.
	Compare() (bool, error)
	CompareUp() (bool, error)
	CompareDown() (bool, error)

	//Attack attacks the entity in front of the turtle.
	Attack() (bool, error)
	AttackUp() (bool, error)
	AttackDown() (bool, error)

	//Suck sucks an item from the inventory in front of the turtle, or from an item floating in the world.
	Suck(count int) (bool, error)
	SuckUp(count int) (bool, error)
	SuckDown(count int) (bool, error)

	//FuelLevel gets the maximum amount of fuel this turtle currently holds.
	FuelLevel() (int, error)
	//Refuel refuels this turtle
	Refuel(count int) (bool, error)
	//FuelLimit gets the maximum amount of fuel this turtle can hold.
	FuelLimit() (int, error)

	//Inspect gets information about the block in front of the turtle.
	Inspect() (bool, Block, error)
	InspectUp() (bool, Block, error)
	InspectDown() (bool, Block, error)

	//Craft crafts a recipe based on the turtle's inventory.
	Craft(limit int) (bool, error)

	io.Closer
	Computer
	Settings
	GPS
}

type turtle struct {
	client connection.Client
}

func NewTurtle(client connection.Client) *turtle {
	t := new(turtle)
	t.client = client
	return t
}

func (t *turtle) Online() bool {
	return t.client.Online()
}

func (t *turtle) Close() error {
	t.Shutdown()
	return t.client.Connection().Close()
}

func (t *turtle) Shutdown() error {
	conn := t.client.Connection()
	return commands.Shutdown(conn.Context(), conn)
}

func (t *turtle) Reboot() error {
	conn := t.client.Connection()
	return commands.Reboot(conn.Context(), conn)
}

func (t *turtle) Version() (string, error) {
	conn := t.client.Connection()
	return commands.Version(conn.Context(), conn)
}

func (t *turtle) ComputerId() (string, error) {
	conn := t.client.Connection()
	return commands.ComputerId(conn.Context(), conn)
}

func (t *turtle) ComputerLabel() (string, error) {
	conn := t.client.Connection()
	return commands.ComputerLabel(conn.Context(), conn)
}

func (t *turtle) SetComputerLabel(label string) error {
	conn := t.client.Connection()
	return commands.SetComputerLabel(conn.Context(), conn, label)
}

func (t *turtle) Uptime() (time.Duration, error) {
	conn := t.client.Connection()
	return commands.Uptime(conn.Context(), conn)
}

func (t *turtle) Time() (float64, error) {
	conn := t.client.Connection()
	return commands.Time(conn.Context(), conn)
}

func (t *turtle) Forward() (bool, error) {
	conn := t.client.Connection()
	return connection.DoActionBool(conn.Context(), conn, "turtle.forward()")
}

func (t *turtle) Back() (bool, error) {
	conn := t.client.Connection()
	return connection.DoActionBool(conn.Context(), conn, "turtle.back()")
}

func (t *turtle) Up() (bool, error) {
	conn := t.client.Connection()
	return connection.DoActionBool(conn.Context(), conn, "turtle.up()")
}

func (t *turtle) Down() (bool, error) {
	conn := t.client.Connection()
	return connection.DoActionBool(conn.Context(), conn, "turtle.down()")
}

func (t *turtle) TurnLeft() (bool, error) {
	conn := t.client.Connection()
	return connection.DoActionBool(conn.Context(), conn, "turtle.turnLeft()")
}

func (t *turtle) TurnRight() (bool, error) {
	conn := t.client.Connection()
	return connection.DoActionBool(conn.Context(), conn, "turtle.turnRight()")
}

func (t *turtle) _detect(command string) bool {
	conn := t.client.Connection()
	res, err := conn.Execute(conn.Context(), command)
	if err != nil {
		return false
	}

	detect, ok := res[0].(bool)
	return ok && detect
}

func (t *turtle) Detect() bool {
	return t._detect("turtle.detect()")
}

func (t *turtle) DetectUp() bool {
	return t._detect("turtle.detectUp()")

}

func (t *turtle) DetectDown() bool {
	return t._detect("turtle.detectDown()")

}

func (t *turtle) Dig() (bool, error) {
	conn := t.client.Connection()
	return connection.DoActionBool(conn.Context(), conn, "turtle.dig()")
}

func (t *turtle) DigDown() (bool, error) {
	conn := t.client.Connection()
	return connection.DoActionBool(conn.Context(), conn, "turtle.digDown()")
}

func (t *turtle) DigUp() (bool, error) {
	conn := t.client.Connection()
	return connection.DoActionBool(conn.Context(), conn, "turtle.digUp()")
}

func (t *turtle) Place() (bool, error) {
	conn := t.client.Connection()
	return connection.DoActionBool(conn.Context(), conn, "turtle.place()")
}

func (t *turtle) PlaceUp() (bool, error) {
	conn := t.client.Connection()
	return connection.DoActionBool(conn.Context(), conn, "turtle.placeUp()")

}

func (t *turtle) PlaceDown() (bool, error) {
	conn := t.client.Connection()
	return connection.DoActionBool(conn.Context(), conn, "turtle.placeDown()")

}

func (t *turtle) Drop(count int) (bool, error) {
	conn := t.client.Connection()
	return connection.DoActionBool(conn.Context(), conn, fmt.Sprintf("turtle.drop(%v)", count))

}

func (t *turtle) DropUp(count int) (bool, error) {
	conn := t.client.Connection()
	return connection.DoActionBool(conn.Context(), conn, fmt.Sprintf("turtle.dropUp(%v)", count))
}

func (t *turtle) DropDown(count int) (bool, error) {
	conn := t.client.Connection()
	return connection.DoActionBool(conn.Context(), conn, fmt.Sprintf("turtle.dropDown(%v)", count))
}

func (t *turtle) Select(slot int) (bool, error) {
	conn := t.client.Connection()
	return connection.DoActionBool(conn.Context(), conn, fmt.Sprintf("turtle.select(%v)", slot))
}

func (t *turtle) SelectedSlot() (int, error) {
	conn := t.client.Connection()
	return connection.DoActionInt(conn.Context(), conn, "turtle.getSelectedSlot()")
}

func (t *turtle) ItemCount(slot int) (int, error) {
	conn := t.client.Connection()
	return connection.DoActionInt(conn.Context(), conn, fmt.Sprintf("turtle.getItemCount(%v)", slot))
}

func (t *turtle) ItemSpace(slot int) (int, error) {
	conn := t.client.Connection()
	return connection.DoActionInt(conn.Context(), conn, fmt.Sprintf("turtle.getItemSpace(%v)", slot))
}

func (t *turtle) ItemDetail(slot int, detailed bool) (map[string]interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (t *turtle) CompareTo(slot int) (bool, error) {
	conn := t.client.Connection()
	return connection.DoActionBool(conn.Context(), conn, fmt.Sprintf("turtle.compareTo(%v)", slot))
}

func (t *turtle) TransferTo(slot, count int) (bool, error) {
	conn := t.client.Connection()
	return connection.DoActionBool(conn.Context(), conn, fmt.Sprintf("turtle.transferTo(%v,%v)", slot, count))
}

func (t *turtle) Compare() (bool, error) {
	conn := t.client.Connection()
	return connection.DoActionBool(conn.Context(), conn, "turtle.compare()")
}

func (t *turtle) CompareUp() (bool, error) {
	conn := t.client.Connection()
	return connection.DoActionBool(conn.Context(), conn, "turtle.compareUp()")

}

func (t *turtle) CompareDown() (bool, error) {
	conn := t.client.Connection()
	return connection.DoActionBool(conn.Context(), conn, "turtle.compareDown()")

}

func (t *turtle) Attack() (bool, error) {
	conn := t.client.Connection()
	return connection.DoActionBool(conn.Context(), conn, "turtle.attack()")

}

func (t *turtle) AttackUp() (bool, error) {
	conn := t.client.Connection()
	return connection.DoActionBool(conn.Context(), conn, "turtle.attackUp()")

}

func (t *turtle) AttackDown() (bool, error) {
	conn := t.client.Connection()
	return connection.DoActionBool(conn.Context(), conn, "turtle.attackDown()")

}

func (t *turtle) Suck(count int) (bool, error) {
	conn := t.client.Connection()
	return connection.DoActionBool(conn.Context(), conn, fmt.Sprintf("turtle.suck(%v)", count))
}

func (t *turtle) SuckUp(count int) (bool, error) {
	conn := t.client.Connection()
	return connection.DoActionBool(conn.Context(), conn, fmt.Sprintf("turtle.suckUp(%v)", count))
}

func (t *turtle) SuckDown(count int) (bool, error) {
	conn := t.client.Connection()
	return connection.DoActionBool(conn.Context(), conn, fmt.Sprintf("turtle.suckDown(%v)", count))
}

func (t *turtle) FuelLevel() (int, error) {
	conn := t.client.Connection()
	return connection.DoActionInt(conn.Context(), conn, "turtle.getFuelLevel()")
}

func (t *turtle) Refuel(count int) (bool, error) {
	conn := t.client.Connection()
	return connection.DoActionBool(conn.Context(), conn, fmt.Sprintf("turtle.refuel(%v)", count))
}

func (t *turtle) FuelLimit() (int, error) {
	conn := t.client.Connection()
	return connection.DoActionInt(conn.Context(), conn, "turtle.getFuelLimit()")
}

func (t *turtle) _doInspect(command string) (bool, Block, error) {
	conn := t.client.Connection()
	res, err := conn.Execute(conn.Context(), command)
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

func (t *turtle) Inspect() (bool, Block, error) {
	return t._doInspect("turtle.inspect()")
}

func (t *turtle) InspectUp() (bool, Block, error) {
	return t._doInspect("turtle.inspectUp()")

}

func (t *turtle) InspectDown() (bool, Block, error) {
	return t._doInspect("turtle.inspectDown()")

}

func (t *turtle) Craft(limit int) (bool, error) {
	conn := t.client.Connection()
	return connection.DoActionBool(conn.Context(), conn, fmt.Sprintf("turtle.craft(%v)", limit))
}

func (t *turtle) _doLocate(timeout time.Duration, debug bool) (int, int, int, error) {
	conn := t.client.Connection()
	res, err := conn.Execute(
		conn.Context(),
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

func (t *turtle) Locate() (int, int, int, error) {
	return t._doLocate(2*time.Second, false)
}

func (t *turtle) LocateWithTimeout(timeout time.Duration) (int, int, int, error) {
	return t._doLocate(timeout, false)
}

func (t *turtle) Define(name string, option ...SettingsOption) error {
	conn := t.client.Connection()
	_, err := conn.Execute(conn.Context(), fmt.Sprintf("settings.define(\"%s\")", name))
	return err
}

func (t *turtle) Undefine(name string) error {
	conn := t.client.Connection()
	_, err := conn.Execute(conn.Context(), fmt.Sprintf("settings.undefine(\"%s\")", name))
	return err
}

func (t *turtle) Set(name, value string) error {
	conn := t.client.Connection()
	_, err := conn.Execute(conn.Context(), fmt.Sprintf("settings.set(\"%s\", \"%s\")", name, value))
	return err
}

func (t *turtle) Unset(name string) error {
	conn := t.client.Connection()
	_, err := conn.Execute(conn.Context(), fmt.Sprintf("settings.unset(\"%s\")", name))
	return err
}

func (t *turtle) Get(name string) (string, error) {
	conn := t.client.Connection()
	res, err := conn.Execute(conn.Context(), fmt.Sprintf("settings.set(\"%s\")", name))
	if err != nil {
		return "", err
	}
	str, ok := res[0].(string)
	if !ok {
		return "", errors.New("not a string")
	}
	return str, nil
}

func (t *turtle) Clear() error {
	conn := t.client.Connection()
	_, err := conn.Execute(conn.Context(), fmt.Sprintf("settings.clear()"))
	return err
}

func (t *turtle) Names() ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (t *turtle) Load(path string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (t *turtle) Save(path string) (bool, error) {
	//TODO implement me
	panic("implement me")
}
