package backend

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	rpctypes "github.com/tharsis/ethermint/rpc/ethereum/types"
)

func (e *EVMBackend) processBlock(block *map[string]interface{}) {
	e.logger.Debug("process block {}", block)
}

func (e *EVMBackend) FeeHistory(blockCount rpc.DecimalOrHex, lastBlock rpc.BlockNumber, rewardPercentiles []float64) (*rpctypes.FeeHistoryResult, error) {
	e.logger.Debug("eth_feeHistory count {}  last block {}  rewardPercentiles {} ", blockCount, lastBlock, rewardPercentiles)

	var blockend int64 = int64(lastBlock)

	if blockend <= 0 {
		blockNumber, err := e.BlockNumber()
		if err != nil {
			return nil, err
		}
		blockend = int64(blockNumber)
	}
	var blockstart int64 = blockend - int64(blockCount)
	var blockcount int64 = blockend - blockstart

	var OldestBlock *hexutil.Big = (*hexutil.Big)(big.NewInt(blockstart))
	var Reward [][]*hexutil.Big = make([][]*hexutil.Big, blockcount)
	var rewardcount = len(rewardPercentiles)
	for i := 0; i < int(blockCount); i++ {
		Reward[i] = make([]*hexutil.Big, rewardcount)
	}

	// fetch block
	for blockid := blockstart; blockid < blockend; blockid++ {
		foundblock, err := e.GetBlockByNumber(rpctypes.BlockNumber(blockid), true)
		if err != nil {
			return nil, err
		}
		e.processBlock(&foundblock)

	}

	var BaseFee []*hexutil.Big = make([]*hexutil.Big, blockcount)
	var GasUsedRatio []float64 = make([]float64, blockcount)

	feeHistory := rpctypes.FeeHistoryResult{
		OldestBlock:  OldestBlock,
		Reward:       Reward,
		BaseFee:      BaseFee,
		GasUsedRatio: GasUsedRatio}

	return &feeHistory, fmt.Errorf("eth_feeHistory not implemented")
}
