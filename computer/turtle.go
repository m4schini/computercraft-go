package computer

import (
	"context"
	"errors"
	"fmt"
	"github.com/m4schini/cc-go/connection"
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

func rpcErr(err error) error {
	return fmt.Errorf("RPC: %v", err)
}

type turtle struct {
	id   string
	conn connection.Connection
}

func NewTurtle(connection connection.Connection) *turtle {
	t := new(turtle)
	t.id = connection.Handshake().Id
	t.conn = connection
	return t
}

func (t *turtle) Close() error {
	t.Shutdown()
	return t.conn.Close()
}

func (t *turtle) _doActionBool(command string) (bool, error) {
	res, err := t.conn.Execute(context.TODO(), command)
	if err != nil {
		return false, rpcErr(err)
	}

	if res != nil && len(res) > 1 {
		err, ok := res[1].(string)
		if ok {
			return false, fmt.Errorf("RPC: %v", err)
		}
	}

	detect, ok := res[0].(bool)
	return ok && detect, nil
}

func (t *turtle) _doActionInt(command string) (int, error) {
	res, err := t.conn.Execute(context.TODO(), command)
	if err != nil {
		return -1, rpcErr(err)
	}

	if res != nil && len(res) > 1 {
		err, ok := res[1].(string)
		if ok {
			return -1, fmt.Errorf("RPC: %v", err)
		}
	}

	num, ok := res[0].(float64)
	if ok {
		return int(num), nil
	} else {
		return -1, errors.New("unexpected datatype")
	}
}

func (t *turtle) UUID() string {
	return t.conn.UUID()
}

func (t *turtle) Shutdown() error {
	_, err := t.conn.Execute(context.TODO(), "os.shutdown()")
	if err != nil {
		return rpcErr(err)
	}
	return nil
}

func (t *turtle) Reboot() error {
	_, err := t.conn.Execute(context.TODO(), "os.reboot()")
	if err != nil {
		return rpcErr(err)
	}
	return nil
}

func (t *turtle) Version() (string, error) {
	res, err := t.conn.Execute(context.TODO(), "os.version()")
	if err != nil {
		return "", rpcErr(err)
	}

	version, ok := res[0].(string)
	if ok {
		return version, nil
	} else {
		return "", errors.New("unexpected datatype")
	}
}

func (t *turtle) ComputerId() (string, error) {
	//TODO fix handshake id to use IP+ID again
	res, err := t.conn.Execute(context.TODO(), "os.getComputerID()")
	if err != nil {
		return "", rpcErr(err)
	}

	if len(res) != 1 {
		return "", errors.New("something went wrong")
	}

	label, ok := res[0].(float64)
	if ok {
		return fmt.Sprintf("%v", label), nil
	}

	return "", errors.New("not the id")
}

func (t *turtle) ComputerLabel() (string, error) {
	res, err := t.conn.Execute(context.TODO(), "os.getComputerLabel()")
	if err != nil {
		return "", rpcErr(err)
	}

	if len(res) != 1 {
		return "", errors.New("something went wrong")
	}

	label, ok := res[0].(string)
	if ok {
		return label, nil
	}

	return "", errors.New("not a label")
}

func (t *turtle) SetComputerLabel(label string) error {
	var err error
	if label == "" {
		_, err = t.conn.Execute(context.TODO(), "os.setComputerLabel()")
	} else {
		_, err = t.conn.Execute(context.TODO(), fmt.Sprintf("os.setComputerLabel(\"%v\")", label))
	}

	if err != nil {
		return rpcErr(err)
	} else {
		return nil
	}
}

func (t *turtle) Uptime() (time.Duration, error) {
	res, err := t.conn.Execute(context.TODO(), "os.clock()")
	if err != nil {
		return 0, err
	}

	uptime, ok := res[0].(float64)
	if !ok {
		return 0, err
	}

	return time.Duration(uptime) * time.Second, nil
}

func (t *turtle) Time() (float64, error) {
	res, err := t.conn.Execute(context.TODO(), "os.time()")
	if err != nil {
		return 0, err
	}

	tme, ok := res[0].(float64)
	if !ok {
		return 0, err
	}

	return tme, nil
}

func (t *turtle) Forward() (bool, error) {
	return t._doActionBool("turtle.forward()")
}

func (t *turtle) Back() (bool, error) {
	return t._doActionBool("turtle.back()")
}

func (t *turtle) Up() (bool, error) {
	return t._doActionBool("turtle.up()")
}

func (t *turtle) Down() (bool, error) {
	return t._doActionBool("turtle.down()")
}

func (t *turtle) TurnLeft() (bool, error) {
	return t._doActionBool("turtle.turnLeft()")
}

func (t *turtle) TurnRight() (bool, error) {
	return t._doActionBool("turtle.turnRight()")
}

func (t *turtle) _detect(command string) bool {
	res, err := t.conn.Execute(context.TODO(), command)
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
	return t._doActionBool("turtle.dig()")
}

func (t *turtle) DigDown() (bool, error) {
	return t._doActionBool("turtle.digDown()")
}

func (t *turtle) DigUp() (bool, error) {
	return t._doActionBool("turtle.digUp()")
}

func (t *turtle) Place() (bool, error) {
	return t._doActionBool("turtle.place()")
}

func (t *turtle) PlaceUp() (bool, error) {
	return t._doActionBool("turtle.placeUp()")

}

func (t *turtle) PlaceDown() (bool, error) {
	return t._doActionBool("turtle.placeDown()")

}

func (t *turtle) Drop(count int) (bool, error) {
	return t._doActionBool(fmt.Sprintf("turtle.drop(%v)", count))

}

func (t *turtle) DropUp(count int) (bool, error) {
	return t._doActionBool(fmt.Sprintf("turtle.dropUp(%v)", count))
}

func (t *turtle) DropDown(count int) (bool, error) {
	return t._doActionBool(fmt.Sprintf("turtle.dropDown(%v)", count))
}

func (t *turtle) Select(slot int) (bool, error) {
	return t._doActionBool(fmt.Sprintf("turtle.select(%v)", slot))
}

func (t *turtle) SelectedSlot() (int, error) {
	return t._doActionInt("turtle.getSelectedSlot()")
}

func (t *turtle) ItemCount(slot int) (int, error) {
	return t._doActionInt(fmt.Sprintf("turtle.getItemCount(%v)", slot))
}

func (t *turtle) ItemSpace(slot int) (int, error) {
	return t._doActionInt(fmt.Sprintf("turtle.getItemSpace(%v)", slot))
}

func (t *turtle) ItemDetail(slot int, detailed bool) (map[string]interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (t *turtle) CompareTo(slot int) (bool, error) {
	return t._doActionBool(fmt.Sprintf("turtle.compareTo(%v)", slot))
}

func (t *turtle) TransferTo(slot, count int) (bool, error) {
	return t._doActionBool(fmt.Sprintf("turtle.transferTo(%v,%v)", slot, count))
}

func (t *turtle) Compare() (bool, error) {
	return t._doActionBool("turtle.compare()")
}

func (t *turtle) CompareUp() (bool, error) {
	return t._doActionBool("turtle.compareUp()")

}

func (t *turtle) CompareDown() (bool, error) {
	return t._doActionBool("turtle.compareDown()")

}

func (t *turtle) Attack() (bool, error) {
	return t._doActionBool("turtle.attack()")

}

func (t *turtle) AttackUp() (bool, error) {
	return t._doActionBool("turtle.attackUp()")

}

func (t *turtle) AttackDown() (bool, error) {
	return t._doActionBool("turtle.attackDown()")

}

func (t *turtle) Suck(count int) (bool, error) {
	return t._doActionBool(fmt.Sprintf("turtle.suck(%v)", count))
}

func (t *turtle) SuckUp(count int) (bool, error) {
	return t._doActionBool(fmt.Sprintf("turtle.suckUp(%v)", count))
}

func (t *turtle) SuckDown(count int) (bool, error) {
	return t._doActionBool(fmt.Sprintf("turtle.suckDown(%v)", count))
}

func (t *turtle) FuelLevel() (int, error) {
	return t._doActionInt("turtle.getFuelLevel()")
}

func (t *turtle) Refuel(count int) (bool, error) {
	return t._doActionBool(fmt.Sprintf("turtle.refuel(%v)", count))
}

func (t *turtle) FuelLimit() (int, error) {
	return t._doActionInt("turtle.getFuelLimit()")
}

func (t *turtle) _doInspect(command string) (bool, Block, error) {
	res, err := t.conn.Execute(context.TODO(), command)
	if err != nil {
		return false, nil, rpcErr(err)
	}

	if len(res) < 2 {
		return false, nil, rpcErr(errors.New("not enough parameter"))
	}

	errMsg, isError := res[1].(string)
	if isError {
		return false, nil, rpcErr(errors.New(errMsg))
	}

	detectedBlock, ok := res[0].(bool)
	if !detectedBlock {
		return false, nil, nil
	}

	data, ok := res[1].(map[string]interface{})
	if !ok {
		return false, nil, rpcErr(errors.New("unexpected datatype"))
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
	return t._doActionBool(fmt.Sprintf("turtle.craft(%v)", limit))
}

func (t *turtle) _doLocate(timeout time.Duration, debug bool) (int, int, int, error) {
	res, err := t.conn.Execute(
		context.TODO(),
		fmt.Sprintf("gps.locate(%v, %v)", int(timeout.Seconds()), debug),
	)
	if err != nil {
		return 0, 0, 0, rpcErr(err)
	}

	if len(res) == 0 {
		return 0, 0, 0, errors.New("position could not be established")
	}

	if len(res) < 3 {
		return 0, 0, 0, errors.New("position could not be established")
	}

	x, ok := res[0].(float64)
	if !ok {
		return 0, 0, 0, errors.New("unexpected datatype")
	}

	y, ok := res[1].(float64)
	if !ok {
		return 0, 0, 0, errors.New("unexpected datatype")
	}

	z, ok := res[2].(float64)
	if !ok {
		return 0, 0, 0, errors.New("unexpected datatype")
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
	_, err := t.conn.Execute(context.TODO(), fmt.Sprintf("settings.define(\"%s\")", name))
	return err
}

func (t *turtle) Undefine(name string) error {
	_, err := t.conn.Execute(context.TODO(), fmt.Sprintf("settings.undefine(\"%s\")", name))
	return err
}

func (t *turtle) Set(name, value string) error {
	_, err := t.conn.Execute(context.TODO(), fmt.Sprintf("settings.set(\"%s\", \"%s\")", name, value))
	return err
}

func (t *turtle) Unset(name string) error {
	_, err := t.conn.Execute(context.TODO(), fmt.Sprintf("settings.unset(\"%s\")", name))
	return err
}

func (t *turtle) Get(name string) (string, error) {
	res, err := t.conn.Execute(context.TODO(), fmt.Sprintf("settings.set(\"%s\")", name))
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
	_, err := t.conn.Execute(context.TODO(), fmt.Sprintf("settings.clear()"))
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
