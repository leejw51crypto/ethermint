package backend

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
)

type OracleBackend struct {
	evm *EVMBackend
}

func (e *OracleBackend) ChainConfig() *params.ChainConfig {
	return nil
}

func (e *OracleBackend) HeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Header, error) {
	return nil, nil

}
func (e *OracleBackend) BlockByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Block, error) {
	return nil, nil
}
func (e *OracleBackend) GetReceipts(ctx context.Context, hash common.Hash) (types.Receipts, error) {
	return nil, nil

}
func (e *OracleBackend) PendingBlockAndReceipts() (*types.Block, types.Receipts) {
	return nil, nil

}
func (e *OracleBackend) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	return nil
}
