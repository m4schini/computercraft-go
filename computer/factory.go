package computer

import (
	"cc-go/connection"
)

func MakeTurtle(conn connection.Connection) *turtle {
	t := NewTurtle(conn)
	return t
}
