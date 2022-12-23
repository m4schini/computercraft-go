package connection

import (
	"fmt"
	"github.com/pkg/errors"
)

func RpcError(err error) error {
	return errors.Wrap(err, "WS-RPC")
}

var UnexpectedDatatypeErr = fmt.Errorf("unexpected datatype")