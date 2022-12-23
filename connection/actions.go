package connection

import (
	"context"
	"errors"
)

func DoActionBool(ctx context.Context, conn Connection, command string) (bool, error) {
	res, err := conn.Execute(ctx, command)
	if err != nil {
		return false, RpcError(err)
	}

	if res != nil && len(res) > 1 {
		err, ok := res[1].(string)
		if ok {
			return false, RpcError(errors.New(err))
		}
	}

	detect, ok := res[0].(bool)
	return ok && detect, nil
}

func DoActionInt(ctx context.Context, conn Connection, command string) (int, error) {
	res, err := conn.Execute(ctx, command)
	if err != nil {
		return -1, RpcError(err)
	}

	if res != nil && len(res) > 1 {
		err, ok := res[1].(string)
		if ok {
			return -1, RpcError(errors.New(err))
		}
	}

	num, ok := res[0].(float64)
	if ok {
		return int(num), nil
	} else {
		return -1, errors.New("unexpected datatype")
	}
}