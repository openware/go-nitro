package rpcproto

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/stretchr/testify/assert"
)

func TestCreateObjectiveRequest(t *testing.T) {
	m := map[string]interface{}{
		"app_data":           []byte("abc"),
		"app_definition":     "0x0000000000000000000123000000000000000000",
		"challenge_duration": uint32(600),
		"counter_party":      "0x111A00868581f73AB42FEEF67D235Ca09ca1E8db",
		"nonce":              uint64(18057610176738954346),
		"outcome": []map[string]interface{}{
			{
				"allocations": []map[string]interface{}{
					{
						"allocation_type": 0,
						"amount":          "100",
						"destination":     "0x000000000000000000000000aaa6628ec44a8a742987ef3a114ddfe2d4f7adce",
						"metadata":        nil,
					},
					{
						"allocation_type": 0,
						"amount":          "100",
						"destination":     "0x000000000000000000000000111a00868581f73ab42feef67d235ca09ca1e8db",
						"metadata":        nil,
					},
				},
				"asset":    "0x0000000000000000000000000000000000000000",
				"metadata": nil,
			},
		},
	}

	r := CreateDirectFundObjectiveRequest(m)

	assert.Equal(t, &directfund.ObjectiveRequest{
		CounterParty:      common.HexToAddress("0x111A00868581f73AB42FEEF67D235Ca09ca1E8db"),
		ChallengeDuration: uint32(600),
		// Outcome:           ,
		AppDefinition: common.HexToAddress("0x0000000000000000000123000000000000000000"),
		AppData:       []byte("abc"),
		Nonce:         uint64(18057610176738954346),
	}, r)
}
