package connection

import (
	"context"
)

type nop struct {
}

func NewNopConnection() (conn *nop) {
	return new(nop)
}

func (n *nop) Execute(ctx context.Context, command string) ([]interface{}, error) {
	return []interface{}{0}, nil
}
