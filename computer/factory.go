package computer

import (
	"github.com/m4schini/computercraft-go/connection"
)

func MakeTurtle(conn connection.Connection) *turtle {
	t := NewTurtle(conn)
	return t
}
