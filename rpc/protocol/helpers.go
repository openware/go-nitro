package rpcproto

import (
	"fmt"
	"reflect"
)

func I2Uint32(v interface{}) uint32 {
	switch v.(type) {
	case float64:
		return uint32(v.(float64))
	case uint32:
		return v.(uint32)
	}
	panic(fmt.Sprintf("invalid type %s", reflect.TypeOf(v)))
}

func I2Uint64(v interface{}) uint64 {
	switch v.(type) {
	case float64:
		return uint64(v.(float64))
	case uint64:
		return v.(uint64)
	}
	panic(fmt.Sprintf("invalid type %s", reflect.TypeOf(v)))
}
