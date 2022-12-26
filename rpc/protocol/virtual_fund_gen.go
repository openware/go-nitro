package rpcproto

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *VirtualFundRequest) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "intermediaries":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				err = msgp.WrapError(err, "Intermediaries")
				return
			}
			if cap(z.Intermediaries) >= int(zb0002) {
				z.Intermediaries = (z.Intermediaries)[:zb0002]
			} else {
				z.Intermediaries = make([]string, zb0002)
			}
			for za0001 := range z.Intermediaries {
				z.Intermediaries[za0001], err = dc.ReadString()
				if err != nil {
					err = msgp.WrapError(err, "Intermediaries", za0001)
					return
				}
			}
		case "counter_party":
			z.CounterParty, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "CounterParty")
				return
			}
		case "challenge_duration":
			z.ChallengeDuration, err = dc.ReadUint32()
			if err != nil {
				err = msgp.WrapError(err, "ChallengeDuration")
				return
			}
		case "outcome":
			var zb0003 uint32
			zb0003, err = dc.ReadArrayHeader()
			if err != nil {
				err = msgp.WrapError(err, "Outcome")
				return
			}
			if cap(z.Outcome) >= int(zb0003) {
				z.Outcome = (z.Outcome)[:zb0003]
			} else {
				z.Outcome = make([]SingleAssetExit, zb0003)
			}
			for za0002 := range z.Outcome {
				err = z.Outcome[za0002].DecodeMsg(dc)
				if err != nil {
					err = msgp.WrapError(err, "Outcome", za0002)
					return
				}
			}
		case "nonce":
			z.Nonce, err = dc.ReadUint64()
			if err != nil {
				err = msgp.WrapError(err, "Nonce")
				return
			}
		case "app_definition":
			z.AppDefinition, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "AppDefinition")
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *VirtualFundRequest) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 6
	// write "intermediaries"
	err = en.Append(0x86, 0xae, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x72, 0x69, 0x65, 0x73)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Intermediaries)))
	if err != nil {
		err = msgp.WrapError(err, "Intermediaries")
		return
	}
	for za0001 := range z.Intermediaries {
		err = en.WriteString(z.Intermediaries[za0001])
		if err != nil {
			err = msgp.WrapError(err, "Intermediaries", za0001)
			return
		}
	}
	// write "counter_party"
	err = en.Append(0xad, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x5f, 0x70, 0x61, 0x72, 0x74, 0x79)
	if err != nil {
		return
	}
	err = en.WriteString(z.CounterParty)
	if err != nil {
		err = msgp.WrapError(err, "CounterParty")
		return
	}
	// write "challenge_duration"
	err = en.Append(0xb2, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x5f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e)
	if err != nil {
		return
	}
	err = en.WriteUint32(z.ChallengeDuration)
	if err != nil {
		err = msgp.WrapError(err, "ChallengeDuration")
		return
	}
	// write "outcome"
	err = en.Append(0xa7, 0x6f, 0x75, 0x74, 0x63, 0x6f, 0x6d, 0x65)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Outcome)))
	if err != nil {
		err = msgp.WrapError(err, "Outcome")
		return
	}
	for za0002 := range z.Outcome {
		err = z.Outcome[za0002].EncodeMsg(en)
		if err != nil {
			err = msgp.WrapError(err, "Outcome", za0002)
			return
		}
	}
	// write "nonce"
	err = en.Append(0xa5, 0x6e, 0x6f, 0x6e, 0x63, 0x65)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.Nonce)
	if err != nil {
		err = msgp.WrapError(err, "Nonce")
		return
	}
	// write "app_definition"
	err = en.Append(0xae, 0x61, 0x70, 0x70, 0x5f, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e)
	if err != nil {
		return
	}
	err = en.WriteString(z.AppDefinition)
	if err != nil {
		err = msgp.WrapError(err, "AppDefinition")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *VirtualFundRequest) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 6
	// string "intermediaries"
	o = append(o, 0x86, 0xae, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x72, 0x69, 0x65, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Intermediaries)))
	for za0001 := range z.Intermediaries {
		o = msgp.AppendString(o, z.Intermediaries[za0001])
	}
	// string "counter_party"
	o = append(o, 0xad, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x5f, 0x70, 0x61, 0x72, 0x74, 0x79)
	o = msgp.AppendString(o, z.CounterParty)
	// string "challenge_duration"
	o = append(o, 0xb2, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x5f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e)
	o = msgp.AppendUint32(o, z.ChallengeDuration)
	// string "outcome"
	o = append(o, 0xa7, 0x6f, 0x75, 0x74, 0x63, 0x6f, 0x6d, 0x65)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Outcome)))
	for za0002 := range z.Outcome {
		o, err = z.Outcome[za0002].MarshalMsg(o)
		if err != nil {
			err = msgp.WrapError(err, "Outcome", za0002)
			return
		}
	}
	// string "nonce"
	o = append(o, 0xa5, 0x6e, 0x6f, 0x6e, 0x63, 0x65)
	o = msgp.AppendUint64(o, z.Nonce)
	// string "app_definition"
	o = append(o, 0xae, 0x61, 0x70, 0x70, 0x5f, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e)
	o = msgp.AppendString(o, z.AppDefinition)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *VirtualFundRequest) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "intermediaries":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Intermediaries")
				return
			}
			if cap(z.Intermediaries) >= int(zb0002) {
				z.Intermediaries = (z.Intermediaries)[:zb0002]
			} else {
				z.Intermediaries = make([]string, zb0002)
			}
			for za0001 := range z.Intermediaries {
				z.Intermediaries[za0001], bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "Intermediaries", za0001)
					return
				}
			}
		case "counter_party":
			z.CounterParty, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "CounterParty")
				return
			}
		case "challenge_duration":
			z.ChallengeDuration, bts, err = msgp.ReadUint32Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "ChallengeDuration")
				return
			}
		case "outcome":
			var zb0003 uint32
			zb0003, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Outcome")
				return
			}
			if cap(z.Outcome) >= int(zb0003) {
				z.Outcome = (z.Outcome)[:zb0003]
			} else {
				z.Outcome = make([]SingleAssetExit, zb0003)
			}
			for za0002 := range z.Outcome {
				bts, err = z.Outcome[za0002].UnmarshalMsg(bts)
				if err != nil {
					err = msgp.WrapError(err, "Outcome", za0002)
					return
				}
			}
		case "nonce":
			z.Nonce, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Nonce")
				return
			}
		case "app_definition":
			z.AppDefinition, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "AppDefinition")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *VirtualFundRequest) Msgsize() (s int) {
	s = 1 + 15 + msgp.ArrayHeaderSize
	for za0001 := range z.Intermediaries {
		s += msgp.StringPrefixSize + len(z.Intermediaries[za0001])
	}
	s += 14 + msgp.StringPrefixSize + len(z.CounterParty) + 19 + msgp.Uint32Size + 8 + msgp.ArrayHeaderSize
	for za0002 := range z.Outcome {
		s += z.Outcome[za0002].Msgsize()
	}
	s += 6 + msgp.Uint64Size + 15 + msgp.StringPrefixSize + len(z.AppDefinition)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *VirtualFundResponse) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "id":
			z.Id, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Id")
				return
			}
		case "channel_id":
			z.ChannelId, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "ChannelId")
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z VirtualFundResponse) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "id"
	err = en.Append(0x82, 0xa2, 0x69, 0x64)
	if err != nil {
		return
	}
	err = en.WriteString(z.Id)
	if err != nil {
		err = msgp.WrapError(err, "Id")
		return
	}
	// write "channel_id"
	err = en.Append(0xaa, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x5f, 0x69, 0x64)
	if err != nil {
		return
	}
	err = en.WriteString(z.ChannelId)
	if err != nil {
		err = msgp.WrapError(err, "ChannelId")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z VirtualFundResponse) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "id"
	o = append(o, 0x82, 0xa2, 0x69, 0x64)
	o = msgp.AppendString(o, z.Id)
	// string "channel_id"
	o = append(o, 0xaa, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x5f, 0x69, 0x64)
	o = msgp.AppendString(o, z.ChannelId)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *VirtualFundResponse) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "id":
			z.Id, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Id")
				return
			}
		case "channel_id":
			z.ChannelId, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "ChannelId")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z VirtualFundResponse) Msgsize() (s int) {
	s = 1 + 3 + msgp.StringPrefixSize + len(z.Id) + 11 + msgp.StringPrefixSize + len(z.ChannelId)
	return
}
