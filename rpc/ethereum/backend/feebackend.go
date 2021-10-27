package backend

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	tmrpctypes "github.com/tendermint/tendermint/rpc/core/types"
	rpctypes "github.com/tharsis/ethermint/rpc/ethereum/types"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"
)

func (e *EVMBackend) processBlock(
	tendermintblock *tmrpctypes.ResultBlock,
	block *map[string]interface{}, rewardPercentiles []float64, onefeehistory *OneFeeHistory) error {

	height := tendermintblock.Block.Height
	e.logger.Debug("processBlock #################")
	e.logger.Debug("height ", height, "   ###################")
	json, jsonerr := json.Marshal(block)
	if jsonerr != nil {
		return jsonerr
	}
	e.logger.Debug(string(json))
	basefee, err := e.BaseFee(height)
	if err != nil {
		return err
	}

	// set basefee
	onefeehistory.BaseFee = basefee

	// set gasused
	onefeehistory.GasUsed = 0.2

	var rewardcount = len(rewardPercentiles)
	onefeehistory.Reward = make([]*big.Int, rewardcount)
	for i := 0; i < rewardcount; i++ {
		onefeehistory.Reward[i] = big.NewInt(2000)
	}

	// check txs
	txs := tendermintblock.Block.Txs
	for _, txBz := range txs {
		tx, err := e.clientCtx.TxConfig.TxDecoder()(txBz)
		if err != nil {
			e.logger.Debug("failed to decode transaction in block", "height", height, "error", err.Error())
			continue
		}

		for _, msg := range tx.GetMsgs() {
			ethMsg, ok := msg.(*evmtypes.MsgEthereumTx)
			if !ok {
				continue
			}

			tx := ethMsg.AsTransaction()
			hash := tx.Hash()
			fmt.Printf("tx=%v hash=%v", tx, hash)

		}
	}

	return nil
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

		tendermintblock, tenderminterr := e.GetTendermintBlockByNumber(rpctypes.BlockNumber(blockid))
		if tenderminterr != nil {
			return nil, err
		}

		var onefeehistory OneFeeHistory = OneFeeHistory{}
		processerr := e.processBlock(tendermintblock, &foundblock, rewardPercentiles, &onefeehistory)
		if processerr != nil {
			return nil, processerr
		}

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
