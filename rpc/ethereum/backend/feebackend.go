package backend

import (
	"fmt"
	"math/big"
	"sort"
	"strconv"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	tmrpctypes "github.com/tendermint/tendermint/rpc/core/types"
	rpctypes "github.com/tharsis/ethermint/rpc/ethereum/types"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"
)

type (
	txGasAndReward struct {
		gasUsed uint64
		reward  *big.Int
	}
	sortGasAndReward []txGasAndReward
)

func (s sortGasAndReward) Len() int { return len(s) }
func (s sortGasAndReward) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s sortGasAndReward) Less(i, j int) bool {
	return s[i].reward.Cmp(s[j].reward) < 0
}

func (e *EVMBackend) processBlock(
	tendermintblock *tmrpctypes.ResultBlock,
	block *map[string]interface{}, rewardPercentiles []float64, blockresult *tmrpctypes.ResultBlockResults, onefeehistory *OneFeeHistory) error {

	height := tendermintblock.Block.Height
	e.logger.Debug("processBlock #################")
	e.logger.Debug("height {} #########", height)
	//json, jsonerr := json.Marshal(block)
	//	if jsonerr != nil {
	//return jsonerr
	//}
	//e.logger.Debug(string(json))
	basefee, err := e.BaseFee(height)
	if err != nil {
		return err
	}

	// set basefee
	onefeehistory.BaseFee = basefee

	// set gasused ratio
	var gasLimit2 = (*block)["gasLimit"].(hexutil.Uint64)
	var gasUsed4 = (*block)["gasUsed"].(*hexutil.Big)
	var gasUsed3 = gasUsed4.ToInt().String()
	//e.logger.Debug("gasUsed3 {}", gasUsed3)
	gasUsed2, _ := strconv.ParseFloat(gasUsed3, 64)
	//e.logger.Debug("gasLimit {}", gasLimit2)
	//e.logger.Debug("gasUsed {}", gasUsed2)
	var gasusedratio float64 = 0
	if gasLimit2 > 0 {
		gasusedratio = float64(gasUsed2) / float64(gasLimit2)
	}
	d1 := fmt.Sprintf("ratio= %f  gasused= %f   gaslimit= %f", gasusedratio, float64(gasUsed2), float64(gasLimit2))
	e.logger.Debug(d1)

	var blockgasused = gasUsed2

	onefeehistory.GasUsed = gasusedratio

	var rewardcount = len(rewardPercentiles)
	onefeehistory.Reward = make([]*big.Int, rewardcount)
	for i := 0; i < rewardcount; i++ {
		onefeehistory.Reward[i] = big.NewInt(2000)
	}

	// check txs
	txs := tendermintblock.Block.Txs
	txresults := blockresult.TxsResults
	txcount := len(txs)
	sorter := make(sortGasAndReward, txcount)

	for i := 0; i < txcount; i++ {
		txBz := txs[i]
		txresult := txresults[i]

		tx, err := e.clientCtx.TxConfig.TxDecoder()(txBz)
		if err != nil {
			e.logger.Debug("failed to decode transaction in block", "height", height, "error", err.Error())
			continue
		}
		gasused := uint64(txresult.GasUsed)
		for _, msg := range tx.GetMsgs() {
			ethMsg, ok := msg.(*evmtypes.MsgEthereumTx)
			if !ok {
				continue
			}
			tx := ethMsg.AsTransaction()
			reward := tx.EffectiveGasTipValue(basefee)
			sorter[i] = txGasAndReward{gasUsed: gasused, reward: reward}
			fmt.Printf("reward %v  gas used %v", reward, gasused)

			hash := tx.Hash()
			fmt.Printf("tx=%v hash=%v", tx, hash)
			break
		}
	}
	sort.Sort(sorter)

	// print
	for index, element := range sorter {
		a := fmt.Sprintf("index= %d  gasused= %v   reward= %v", index, element.gasUsed, element.reward)
		e.logger.Debug(a)
	}

	var txIndex int
	sumGasUsed := uint64(0)
	if len(sorter) > 0 {
		sumGasUsed = sorter[0].gasUsed
	}
	for i, p := range rewardPercentiles {
		thresholdGasUsed := uint64(float64(blockgasused) * p / 100)
		for sumGasUsed < thresholdGasUsed && txIndex < txcount-1 {
			b := fmt.Sprintf("index %d  txindex %d  percent %f  thresholdGasUsed %d sumGasUsed %d", i, txIndex, p, thresholdGasUsed, sumGasUsed)
			e.logger.Debug(b)
			txIndex++
			sumGasUsed += sorter[txIndex].gasUsed
		}
		b := fmt.Sprintf("index %d    percent %f  thresholdGasUsed %d sumGasUsed %d", i, p, thresholdGasUsed, sumGasUsed)
		e.logger.Debug(b)

		chosenreward := big.NewInt(0)
		if 0 <= txIndex && txIndex < len(sorter) {
			chosenreward = sorter[txIndex].reward
			e.logger.Debug(fmt.Sprintf("chosen txindex %d  reward %s", txIndex, chosenreward))
		}
		onefeehistory.Reward[i] = chosenreward
	}

	e.logger.Debug(fmt.Sprintf("######## gasLimit %d", gasLimit2))
	e.logger.Debug(fmt.Sprintf("######## gasUsed %f ", gasUsed2))
	e.logger.Debug(fmt.Sprintf("######## sumGasUsed %d ", sumGasUsed))

	return nil
}

type OneFeeHistory struct {
	BaseFee *big.Int
	Reward  []*big.Int
	GasUsed float64
}

func (e *EVMBackend) FeeHistory(userblockCount rpc.DecimalOrHex, lastBlock rpc.BlockNumber, rewardPercentiles []float64) (*rpctypes.FeeHistoryResult, error) {
	e.logger.Debug("eth_feeHistory count {}   ", userblockCount)

	var blockend int64 = int64(lastBlock)

	if blockend <= 0 {
		blockNumber, err := e.BlockNumber()
		if err != nil {
			return nil, err
		}
		blockend = int64(blockNumber)
	}
	var blockstart int64 = blockend - int64(userblockCount)
	if blockstart < 0 {
		blockstart = 0
	}

	var blockcount int64 = blockend - blockstart

	var OldestBlock *hexutil.Big = (*hexutil.Big)(big.NewInt(blockstart))

	// prepare space
	var Reward [][]*hexutil.Big = make([][]*hexutil.Big, blockcount)
	var rewardcount = len(rewardPercentiles)
	for i := 0; i < int(blockcount); i++ {
		Reward[i] = make([]*hexutil.Big, rewardcount)
	}
	var BaseFee []*hexutil.Big = make([]*hexutil.Big, blockcount)
	var GasUsedRatio []float64 = make([]float64, blockcount)

	// fetch block
	for blockid := blockstart; blockid < blockend; blockid++ {
		index := int32(blockid - blockstart)
		// eth block
		foundblock, err := e.GetBlockByNumber(rpctypes.BlockNumber(blockid), true)
		if err != nil {
			return nil, err
		}

		// tendermint block
		tendermintblock, tenderminterr := e.GetTendermintBlockByNumber(rpctypes.BlockNumber(blockid))
		if tenderminterr != nil {
			return nil, err
		}

		// block result
		foundblockresult, err := e.clientCtx.Client.BlockResults(e.ctx, &tendermintblock.Block.Height)
		if err != nil {
			e.logger.Debug("EthBlockFromTendermint block result not found", "height", tendermintblock.Block.Height, "error", err.Error())
			return nil, err
		}

		var onefeehistory OneFeeHistory = OneFeeHistory{}
		processerr := e.processBlock(tendermintblock, &foundblock, rewardPercentiles, foundblockresult, &onefeehistory)
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
