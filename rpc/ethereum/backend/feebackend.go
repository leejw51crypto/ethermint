package backend

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	rpctypes "github.com/tharsis/ethermint/rpc/ethereum/types"
)

func (e *EVMBackend) processBlock(block *map[string]interface{}, rewardPercentiles []float64, onefeehistory *OneFeeHistory) {
	e.logger.Debug("processBlock#################")
	onefeehistory.BaseFee = big.NewInt(100)
	onefeehistory.GasUsed = 0.2

	var rewardcount = len(rewardPercentiles)
	onefeehistory.Reward = make([]*big.Int, rewardcount)
	for i := 0; i < rewardcount; i++ {
		onefeehistory.Reward[i] = big.NewInt(2000)
	}
}

type OneFeeHistory struct {
	BaseFee *big.Int
	Reward  []*big.Int
	GasUsed float64
}

func (e *EVMBackend) FeeHistory(blockCount rpc.DecimalOrHex, lastBlock rpc.BlockNumber, rewardPercentiles []float64) (*rpctypes.FeeHistoryResult, error) {
	e.logger.Debug("eth_feeHistory count {}   ", blockCount)

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

	// prepare space
	var Reward [][]*hexutil.Big = make([][]*hexutil.Big, blockcount)
	var rewardcount = len(rewardPercentiles)
	for i := 0; i < int(blockCount); i++ {
		Reward[i] = make([]*hexutil.Big, rewardcount)
	}
	var BaseFee []*hexutil.Big = make([]*hexutil.Big, blockcount)
	var GasUsedRatio []float64 = make([]float64, blockcount)

	// fetch block
	for blockid := blockstart; blockid < blockend; blockid++ {
		index := int32(blockid - blockstart)
		foundblock, err := e.GetBlockByNumber(rpctypes.BlockNumber(blockid), true)
		if err != nil {
			return nil, err
		}
		var onefeehistory OneFeeHistory = OneFeeHistory{}
		e.processBlock(&foundblock, rewardPercentiles, &onefeehistory)

		// iterate
		BaseFee[index] = (*hexutil.Big)(onefeehistory.BaseFee)
		GasUsedRatio[index] = onefeehistory.GasUsed
		for j := 0; j < rewardcount; j++ {
			Reward[index][j] = (*hexutil.Big)(onefeehistory.Reward[j])
		}

	}

	feeHistory := rpctypes.FeeHistoryResult{
		OldestBlock:  OldestBlock,
		Reward:       Reward,
		BaseFee:      BaseFee,
		GasUsedRatio: GasUsedRatio}
	return &feeHistory, nil
	//return &feeHistory, fmt.Errorf("eth_feeHistory not implemented")
}
