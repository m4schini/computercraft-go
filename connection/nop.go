package connection

import (
	"context"
)

type nop struct {
}

func (n *nop) RemoteHost() string {
	return ""
}

func (n *nop) Id() string {
	return ""
}

func NewNopConnection() *nop {
	return new(nop)
}

func (n *nop) Execute(ctx context.Context, command string) ([]interface{}, error) {
	return []interface{}{0}, nil
}

func (n *nop) Context() context.Context {
	return context.Background()
}

func (n *nop) Device() DeviceType {
	return ""
}

func (n *nop) Close() error {
	return nil
}
