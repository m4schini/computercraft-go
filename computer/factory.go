package computer

import (
	"github.com/m4schini/cc-go/connection"
)

func MakeTurtle(conn connection.Connection) *turtle {
	t := NewTurtle(conn)
	return t
}
