package computer

import (
	"context"
	"fmt"
	"github.com/m4schini/computercraft-go/connection"
	"strings"
)

const (
	PeripheralModuleName = "peripheral"
)

type PeripheralType string

type ItemDetail map[string]any

type Peripheral interface {
	Names() ([]string, error)
	IsPresent(name string) (bool, error)
	GetType(name string) ([]string, error)
	HasType(name string, peripheralType PeripheralType) (bool, error)
	GetMethods(name string) ([]string, error)
	Call(name, method string, args ...any) ([]any, error)
}

func Names(ctx context.Context, conn connection.Connection) ([]string, error) {
	response, err := conn.Execute(ctx, PeripheralModuleName+".getNames()")
	if err != nil || len(response) != 1 {
		return []string{}, err
	}

	names, ok := response[0].([]string)
	if !ok {
		return []string{}, nil
	}

	return names, nil
}

func IsPresent(ctx context.Context, conn connection.Connection, name string) (bool, error) {
	return connection.DoActionBool(ctx, conn, fmt.Sprintf("%v.isPresent(\"%v\")", PeripheralModuleName, name))
}

func GetType(ctx context.Context, conn connection.Connection, name string) ([]string, error) {
	response, err := conn.Execute(ctx, fmt.Sprintf("%v.getType(\"%v\")", PeripheralModuleName, name))
	if err != nil || len(response) != 1 {
		return []string{}, err
	}

	names, ok := response[0].([]string)
	if !ok {
		return []string{}, nil
	}

	return names, nil
}

func HasType(ctx context.Context, conn connection.Connection, name string, peripheralType PeripheralType) (bool, error) {
	response, err := conn.Execute(ctx, fmt.Sprintf("%v.hasType(\"%v\", \"%v\")", PeripheralModuleName, name, peripheralType))
	if err != nil || len(response) != 1 {
		return false, err
	}

	hasType, ok := response[0].(bool)
	if !ok {
		return false, nil
	}

	return hasType, nil
}

func GetMethods(ctx context.Context, conn connection.Connection, name string) ([]string, error) {
	response, err := conn.Execute(ctx, fmt.Sprintf("%v.getMethods(\"%v\")", PeripheralModuleName, name))
	if err != nil || len(response) != 1 {
		return []string{}, err
	}

	methods, ok := response[0].([]string)
	if !ok {
		return []string{}, nil
	}

	return methods, nil
}

func Call(ctx context.Context, conn connection.Connection, name, method string, args ...any) ([]any, error) {
	arguments := make([]string, len(args))

	for i, arg := range args {
		arguments[i] = fmt.Sprintf("%v", arg)
	}

	response, err := conn.Execute(ctx, fmt.Sprintf("%v.call(\"%v\", \"%v\", \"%v\")",
		PeripheralModuleName,
		name,
		method,
		strings.Join(arguments, "\", \"")))
	if err != nil || len(response) != 1 {
		return []any{}, err
	}

	return response, nil
}
