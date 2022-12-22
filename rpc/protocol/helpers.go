package rpcproto

import (
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/statechannels/go-nitro/types"
)

func I2Uint8(v interface{}) uint8 {
	switch v.(type) {
	case int64:
		return uint8(v.(int64))
	case float64:
		return uint8(v.(float64))
	case uint32:
		return uint8(v.(uint32))
	case uint8:
		return v.(uint8)
	}
	panic(fmt.Sprintf("invalid type %s", reflect.TypeOf(v)))
}

func I2Uint32(v interface{}) uint32 {
	switch v.(type) {
	case int64:
		return uint32(v.(int64))
	case float64:
		return uint32(v.(float64))
	case uint32:
		return v.(uint32)
	}
	panic(fmt.Sprintf("invalid type %s", reflect.TypeOf(v)))
}

func I2Uint64(v interface{}) uint64 {
	switch v.(type) {
	case int64:
		return uint64(v.(int64))
	case float64:
		return uint64(v.(float64))
	case uint64:
		return v.(uint64)
	}
	panic(fmt.Sprintf("invalid type %s", reflect.TypeOf(v)))
}

func I2Uint256(v interface{}) *types.Uint256 {
	switch v.(type) {
	case string:
		bigInt, ok := math.ParseBig256(v.(string))
		if !ok {
			panic(fmt.Sprintf("parsing to bigint failed. val: %s", v.(string)))
		}
		return bigInt
	}
	panic(fmt.Sprintf("invalid type %s", reflect.TypeOf(v)))
}
