package computer

import (
	"context"
	"time"
)

type Side string

const (
	SideTop    Side = "top"
	SideBottom Side = "bottom"
	SideLeft   Side = "left"
	SideRight  Side = "right"
	SideFront  Side = "front"
	SideBack   Side = "back"
)

type DeviceType string

const (
	DeviceComputer DeviceType = "computer"
	DeviceTurtle   DeviceType = "turtle"
	DevicePocket   DeviceType = "pocket"
	DeviceUnknown  DeviceType = "unknown"
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
}

type Turtle interface {
	// Forward moves the turtle forward one block.
	Forward(ctx context.Context) (bool, error)
	// Back moves the turtle backwards one block.
	Back(ctx context.Context) (bool, error)
	// Up moves the turtle up one block.
	Up(ctx context.Context) (bool, error)
	// Down moves the turtle down one block.
	Down(ctx context.Context) (bool, error)
	// TurnLeft rotates the turtle 90 degrees to the left.
	TurnLeft(ctx context.Context) (bool, error)
	// TurnRight rotates the turtle 90 degrees to the right.
	TurnRight(ctx context.Context) (bool, error)

	// Dig attempts to break the block in front of the turtle.
	Dig(ctx context.Context) (bool, error)
	// DigDown attempts to break the block below the turtle.
	DigDown(ctx context.Context) (bool, error)
	// DigUp attempts to break the block above the turtle.
	DigUp(ctx context.Context) (bool, error)

	// Place places a block or item into the world in front of the turtle.
	Place(ctx context.Context) (bool, error)
	// PlaceUp places a block or item into the world above the turtle.
	PlaceUp(ctx context.Context) (bool, error)
	// PlaceDown places a block or item into the world below the turtle.
	PlaceDown(ctx context.Context) (bool, error)

	// Drop drops the currently selected stack into the inventory in front of
	//the turtle, or as an item into the world if there is no inventory.
	Drop(ctx context.Context, count int) (bool, error)
	// DropUp drops the currently selected stack into the inventory above the
	//turtle, or as an item into the world if there is no inventory.
	DropUp(ctx context.Context, count int) (bool, error)
	// DropDown drops the currently selected stack into the inventory in front
	//of the turtle, or as an item into the world if there is no inventory.
	DropDown(ctx context.Context, count int) (bool, error)

	// Select changes the currently selected slot.
	Select(ctx context.Context, slot int) (bool, error)
	// SelectedSlot gets the currently selected slot.
	SelectedSlot(ctx context.Context) (int, error)
	// ItemCount gets the number of items in the given slot.
	ItemCount(ctx context.Context, slot int) (int, error)
	// ItemSpace gets the remaining number of items which may be stored in this stack.
	ItemSpace(ctx context.Context, slot int) (int, error)
	// ItemDetail gets detailed information about the items in the given slot.
	ItemDetail(ctx context.Context, slot int, detailed bool) (map[string]interface{}, error)
	// CompareTo compares the item in the currently selected slot to the item in another slot.
	CompareTo(ctx context.Context, slot int) (bool, error)
	// TransferTo moves an item from the selected slot to another one.
	TransferTo(ctx context.Context, slot, count int) (bool, error)

	// Detect checks if there is a solid block in front of the turtle. In this case,
	//solid refers to any non-air or liquid block.
	Detect(ctx context.Context) bool
	// DetectUp checks if there is a solid block above the turtle. In this case,
	//solid refers to any non-air or liquid block.
	DetectUp(ctx context.Context) bool
	// DetectDown checks if there is a solid block below the turtle. In this case,
	//solid refers to any non-air or liquid block.
	DetectDown(ctx context.Context) bool

	// Compare checks if the block in front of the turtle is equal to the item in the currently selected slot.
	Compare(ctx context.Context) (bool, error)
	// CompareUp checks if the block above the turtle is equal to the item in the currently selected slot.
	CompareUp(ctx context.Context) (bool, error)
	// CompareDown checks if the block below the turtle is equal to the item in the currently selected slot.
	CompareDown(ctx context.Context) (bool, error)

	// Attack attacks the entity in front of the turtle.
	Attack(ctx context.Context) (bool, error)
	// AttackUp attacks the entity above the turtle.
	AttackUp(ctx context.Context) (bool, error)
	// AttackDown attacks the entity below the turtle.
	AttackDown(ctx context.Context) (bool, error)

	// Suck sucks an item from the inventory in front of the turtle, or from an item floating in the world.
	Suck(ctx context.Context, count int) (bool, error)
	// SuckUp sucks an item from the inventory above the turtle, or from an item floating in the world.
	SuckUp(ctx context.Context, count int) (bool, error)
	// SuckDown sucks an item from the inventory below the turtle, or from an item floating in the world.
	SuckDown(ctx context.Context, count int) (bool, error)

	// FuelLevel gets the maximum amount of fuel this turtle currently holds.
	FuelLevel(ctx context.Context) (int, error)
	// Refuel refuels this turtle
	Refuel(ctx context.Context, count int) (bool, error)
	// FuelLimit gets the maximum amount of fuel this turtle can hold.
	FuelLimit(ctx context.Context) (int, error)

	// Inspect gets information about the block in front of the turtle.
	Inspect(ctx context.Context) (bool, Block, error)
	// InspectUp gets information about the block above the turtle.
	InspectUp(ctx context.Context) (bool, Block, error)
	// InspectDown gets information about the block below the turtle.
	InspectDown(ctx context.Context) (bool, Block, error)

	// Craft crafts a recipe based on the turtle's inventory.
	Craft(ctx context.Context, limit int) (bool, error)
}
