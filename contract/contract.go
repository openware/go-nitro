package contract

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type StateChannelContract interface {
	Transfer(opts *bind.TransactOpts, assetIndex *big.Int, fromChannelId [32]byte, outcomeBytes []byte, stateHash [32]byte, indices []*big.Int) (*types.Transaction, error)
	TransferAllAssets(opts *bind.TransactOpts, channelId [32]byte, outcomeBytes []byte, stateHash [32]byte) (*types.Transaction, error)
	Deposit(opts *bind.TransactOpts, asset common.Address, channelId [32]byte, expectedHeld *big.Int, amount *big.Int) (*types.Transaction, error)
	ValidTransition(opts *bind.CallOpts, nParticipants *big.Int, isFinalAB [2]bool, ab [2]IForceMoveAppVariablePart, turnNumB *big.Int, appDefinition common.Address) (bool, error)
	GetChainID(opts *bind.CallOpts) (*big.Int, error)
}

type Client struct {
	Contract   StateChannelContract
	EthNetwork *ethclient.Client
}

func NewClient(contractAddr, rpcUrl string) (*Client, error) {
	contractAddress := common.HexToAddress(contractAddr)
	ethNetwork, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, err
	}

	adjucator, err := NewNitroAdjucator(contractAddress, ethNetwork)
	if err != nil {
		return nil, err
	}

	return &Client{
		Contract:   adjucator,
		EthNetwork: ethNetwork,
	}, nil
}
