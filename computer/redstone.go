package computer

type Side string

const (
	SideTop    Side = "top"
	SideBottom Side = "bottom"
	SideLeft   Side = "left"
	SideRight  Side = "right"
	SideFront  Side = "front"
	SideBack   Side = "back"
)

type Redstoner interface {
	//SetOutput turns the redstone signal of a specific side on or off.
	SetOutput(side Side, on bool) error
	//SetAnalogOutput sets the redstone signal strength for a specific side.
	SetAnalogOutput(side Side, value int) error
	//Output gets the current redstone output of a specific side.
	Output(side Side) (bool, int, error)
	//Input gets the current redstone input of a specific side.
	Input(side Side) error
}
