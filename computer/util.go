package computer

import (
	"context"
	"fmt"
	"github.com/m4schini/computercraft-go/connection"
)

type ModuleNotPresentError string

func (m ModuleNotPresentError) Error() string {
	return fmt.Sprintf("module (%v) not present", m)
}

func HasModule(conn connection.Connection, ctx context.Context, moduleName string) error {
	response, err := conn.Execute(ctx, fmt.Sprintf("%v ~= nil", moduleName))
	if err != nil || len(response) != 1 {
		return ModuleNotPresentError(moduleName)
	}

	hasModule, ok := response[0].(bool)
	if !ok || !hasModule {
		return ModuleNotPresentError(moduleName)
	}

	return nil
}

func IsTurtle(ctx context.Context, conn connection.Connection) (bool, error) {
	return connection.DoActionBool(ctx, conn, "turtle ~= nil")
}

func IsPocket(ctx context.Context, conn connection.Connection) (bool, error) {
	return connection.DoActionBool(ctx, conn, "pocket ~= nil")
}

func GetDeviceType(ctx context.Context, conn connection.Connection) (DeviceType, error) {
	isTurtle, err := IsTurtle(ctx, conn)
	if err != nil {
		return DeviceUnknown, err
	}
	if isTurtle {
		return DeviceTurtle, nil
	}

	isPocket, err := IsPocket(ctx, conn)
	if err != nil {
		return DeviceUnknown, err
	}
	if isPocket {
		return DevicePocket, nil
	}

	return DeviceComputer, nil
}
