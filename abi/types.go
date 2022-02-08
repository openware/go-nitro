package abi

import "github.com/ethereum/go-ethereum/accounts/abi"

// To encode objects as bytes, we need to construct an encoder, using abi.Arguments.
// An instance of abi.Arguments implements two functions relevant to us:
// - `Pack`, which packs go values for a given struct into bytes.
// - `unPack`, which unpacks bytes into go values
// To construct an abi.Arguments instance, we need to supply an array of "types", which are
// actually go values. The following types are used when encoding a state

// Uint256 is the Uint256 type for abi encoding
var Uint256, _ = abi.NewType("uint256", "uint256", nil)

// Uint48 is the Uint48 type for abi encoding
var Uint48, _ = abi.NewType("uint48", "uint48", nil)

// Bool is the bool type for abi encoding
var Bool, _ = abi.NewType("bool", "bool", nil)

// Destination is the bytes32 type for abi encoding
var Destination, _ = abi.NewType("bytes32", "address", nil)

// Bytes is the bytes type for abi encoding
var Bytes, _ = abi.NewType("bytes", "bytes", nil)

// Bytes32 is the bytes32 type for abi encoding
var Bytes32, _ = abi.NewType("bytes32", "bytes32", nil)

// AddressArray is the address[] type for abi encoding
var AddressArray, _ = abi.NewType("address[]", "address[]", nil)

// Address is the Address type for abi encoding
var Address, _ = abi.NewType("address", "address", nil)
