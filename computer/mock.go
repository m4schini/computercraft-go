package computer

import (
	"fmt"
	"time"
)

type mockTurtle struct {
	log func(str string)
}

func NewMockTurtle(log func(str string)) *mockTurtle {
	t := new(mockTurtle)
	t.log = log
	return t
}

func (t *mockTurtle) Forward() (bool, error) {
	t.log("turtle.forward()")
	return true, nil
}

func (t *mockTurtle) Back() (bool, error) {
	t.log("turtle.back()")
	return true, nil
}

func (t *mockTurtle) Up() (bool, error) {
	t.log("turtle.up()")
	return true, nil
}

func (t *mockTurtle) Down() (bool, error) {
	t.log("turtle.down()")
	return true, nil
}

func (t *mockTurtle) TurnLeft() (bool, error) {
	t.log("turtle.turnLeft()")
	return true, nil
}

func (t *mockTurtle) TurnRight() (bool, error) {
	t.log("turtle.turnRight()")
	return true, nil
}

func (t *mockTurtle) Dig() (bool, error) {
	t.log("turtle.dig()")
	return true, nil
}

func (t *mockTurtle) DigDown() (bool, error) {
	t.log("turtle.digDown()")
	return true, nil
}

func (t *mockTurtle) DigUp() (bool, error) {
	t.log("turtle.digUp()")
	return true, nil
}

func (t *mockTurtle) Place() (bool, error) {
	t.log("turtle.place()")
	return true, nil
}

func (t *mockTurtle) PlaceUp() (bool, error) {
	t.log("turtle.placeUp()")
	return true, nil
}

func (t *mockTurtle) PlaceDown() (bool, error) {
	t.log("turtle.placeDown()")
	return true, nil
}

func (t *mockTurtle) Drop(count int) (bool, error) {
	t.log(fmt.Sprintf("turtle.drop(%v)", count))
	return true, nil
}

func (t *mockTurtle) DropUp(count int) (bool, error) {
	t.log(fmt.Sprintf("turtle.dropUp(%v)", count))
	return true, nil
}

func (t *mockTurtle) DropDown(count int) (bool, error) {
	t.log(fmt.Sprintf("turtle.dropDown(%v)", count))
	return true, nil
}

func (t *mockTurtle) Select(slot int) (bool, error) {
	t.log(fmt.Sprintf("turtle.select(%v)", slot))
	return true, nil
}

func (t *mockTurtle) SelectedSlot() (int, error) {
	t.log(fmt.Sprintf("turtle.selectedSlot()"))
	return 0, nil
}

func (t *mockTurtle) ItemCount(slot int) (int, error) {
	t.log(fmt.Sprintf("turtle.itemCount(%v)", slot))
	return 0, nil
}

func (t *mockTurtle) ItemSpace(slot int) (int, error) {
	t.log(fmt.Sprintf("turtle.itemSpace(%v)", slot))
	return 64, nil
}

func (t *mockTurtle) ItemDetail(slot int, detailed bool) (map[string]interface{}, error) {
	t.log(fmt.Sprintf("turtle.itemDetail(%v, %v)", slot, detailed))
	if detailed {
		return nil, nil
	} else {
		return nil, nil
	}
}

func (t *mockTurtle) CompareTo(slot int) (bool, error) {
	t.log(fmt.Sprintf("turtle.compareTo(%v)", slot))
	return true, nil
}

func (t *mockTurtle) TransferTo(slot, count int) (bool, error) {
	t.log(fmt.Sprintf("turtle.transferTo(%v, %v)", slot, count))
	return true, nil
}

func (t *mockTurtle) Detect() bool {
	t.log("turtle.detect()")
	return true
}

func (t *mockTurtle) DetectUp() bool {
	t.log("turtle.detectUp()")
	return true
}

func (t *mockTurtle) DetectDown() bool {
	t.log("turtle.detectDown()")
	return true
}

func (t *mockTurtle) Compare() (bool, error) {
	t.log("turtle.compare()")
	return true, nil
}

func (t *mockTurtle) CompareUp() (bool, error) {
	t.log("turtle.compareUp()")
	return true, nil
}

func (t *mockTurtle) CompareDown() (bool, error) {
	t.log("turtle.compareDown()")
	return true, nil
}

func (t *mockTurtle) Attack() (bool, error) {
	t.log("turtle.attack()")
	return true, nil
}

func (t *mockTurtle) AttackUp() (bool, error) {
	t.log("turtle.attackUp()")
	return true, nil
}

func (t *mockTurtle) AttackDown() (bool, error) {
	t.log("turtle.attackDown()")
	return true, nil
}

func (t *mockTurtle) Suck(count int) (bool, error) {
	t.log(fmt.Sprintf("turtle.suck(%v)", count))
	return true, nil
}

func (t *mockTurtle) SuckUp(count int) (bool, error) {
	t.log(fmt.Sprintf("turtle.suckUp(%v)", count))
	return true, nil
}

func (t *mockTurtle) SuckDown(count int) (bool, error) {
	t.log(fmt.Sprintf("turtle.suckDown(%v)", count))
	return true, nil
}

func (t *mockTurtle) FuelLevel() (int, error) {
	t.log(fmt.Sprintf("turtle.fuelLevel()"))
	return 20000, nil
}

func (t *mockTurtle) Refuel(count int) (bool, error) {
	t.log(fmt.Sprintf("turtle.refuel(%v)", count))
	return true, nil
}

func (t *mockTurtle) FuelLimit() (int, error) {
	t.log(fmt.Sprintf("turtle.suck()"))
	return 20000, nil
}

func (t *mockTurtle) Inspect() (bool, Block, error) {
	t.log("turtle.inspect()")
	return false, nil, nil
}

func (t *mockTurtle) InspectUp() (bool, Block, error) {
	t.log("turtle.inspectUp()")
	return false, nil, nil
}

func (t *mockTurtle) InspectDown() (bool, Block, error) {
	t.log("turtle.inspectDown()")
	return false, nil, nil
}

func (t *mockTurtle) Craft(limit int) (bool, error) {
	t.log(fmt.Sprintf("turtle.craft(%v)", limit))
	return true, nil
}

func (t *mockTurtle) Close() error {
	t.log(fmt.Sprintf("closed"))
	t.log = func(str string) {

	}
	return nil
}

func (t *mockTurtle) Shutdown() error {
	t.log("os.shutdown()")
	return nil
}

func (t *mockTurtle) Reboot() error {
	t.log("os.reboot()")
	return nil
}

func (t *mockTurtle) Version() (string, error) {
	t.log("os.shutdown()")
	return "ComputerCraft OS 1.18", nil
}

func (t *mockTurtle) ComputerId() (string, error) {
	t.log("os.computerId()")
	return "0.0.0.0#1", nil
}

func (t *mockTurtle) ComputerLabel() (string, error) {
	t.log("os.computerLabel()")
	return "label", nil
}

func (t *mockTurtle) SetComputerLabel(label string) error {
	//TODO implement me
	panic("implement me")
}

func (t *mockTurtle) Uptime() (time.Duration, error) {
	//TODO implement me
	panic("implement me")
}

func (t *mockTurtle) Time() (float64, error) {
	//TODO implement me
	panic("implement me")
}

func (t *mockTurtle) Define(name string, option ...SettingsOption) error {
	//TODO implement me
	panic("implement me")
}

func (t *mockTurtle) Undefine(name string) error {
	//TODO implement me
	panic("implement me")
}

func (t *mockTurtle) Set(name, value string) error {
	//TODO implement me
	panic("implement me")
}

func (t *mockTurtle) Unset(name string) error {
	//TODO implement me
	panic("implement me")
}

func (t *mockTurtle) Get(name string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (t *mockTurtle) Clear() error {
	//TODO implement me
	panic("implement me")
}

func (t *mockTurtle) Names() ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (t *mockTurtle) Load(path string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (t *mockTurtle) Save(path string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (t *mockTurtle) Locate() (int, int, int, error) {
	//TODO implement me
	panic("implement me")
}

func (t *mockTurtle) LocateWithTimeout(timeout time.Duration) (int, int, int, error) {
	//TODO implement me
	panic("implement me")
}
