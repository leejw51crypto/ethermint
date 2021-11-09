package backend

import (
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
	tendermintBlock *tmrpctypes.ResultBlock,
	ethBlock *map[string]interface{}, rewardPercentiles []float64, tendermintBlockResult *tmrpctypes.ResultBlockResults, targetOneFeeHistory *OneFeeHistory) error {

	blockHeight := tendermintBlock.Block.Height
	blockBaseFee, err := e.BaseFee(blockHeight)
	if err != nil {
		return err
	}

	// set basefee
	targetOneFeeHistory.BaseFee = blockBaseFee

	// set gasused ratio
	var gasLimitUint64 = (*ethBlock)["gasLimit"].(hexutil.Uint64)
	var gasUsedBig = (*ethBlock)["gasUsed"].(*hexutil.Big)
	var gasUsedString = gasUsedBig.ToInt().String()
	gasusedfloat, _ := strconv.ParseFloat(gasUsedString, 64)
	var gasUsedRatio float64 = 0
	if gasLimitUint64 > 0 {
		gasUsedRatio = float64(gasusedfloat) / float64(gasLimitUint64)
	}
	var blockGasUsed = gasusedfloat
	targetOneFeeHistory.GasUsed = gasUsedRatio

	var rewardCount = len(rewardPercentiles)
	targetOneFeeHistory.Reward = make([]*big.Int, rewardCount)
	for i := 0; i < rewardCount; i++ {
		targetOneFeeHistory.Reward[i] = big.NewInt(2000)
	}

	// check tendermintTxs
	tendermintTxs := tendermintBlock.Block.Txs
	tendermintTxResults := tendermintBlockResult.TxsResults
	tendermintTxCount := len(tendermintTxs)
	sorter := make(sortGasAndReward, tendermintTxCount)

	for i := 0; i < tendermintTxCount; i++ {
		eachTendermintTx := tendermintTxs[i]
		eachTendermintTxResult := tendermintTxResults[i]

		tx, err := e.clientCtx.TxConfig.TxDecoder()(eachTendermintTx)
		if err != nil {
			e.logger.Debug("failed to decode transaction in block", "height", blockHeight, "error", err.Error())
			continue
		}
		txGasUsed := uint64(eachTendermintTxResult.GasUsed)
		for _, msg := range tx.GetMsgs() {
			ethMsg, ok := msg.(*evmtypes.MsgEthereumTx)
			if !ok {
				continue
			}
			tx := ethMsg.AsTransaction()
			reward := tx.EffectiveGasTipValue(blockBaseFee)
			sorter[i] = txGasAndReward{gasUsed: txGasUsed, reward: reward}
			break
		}
	}
	sort.Sort(sorter)

	var txIndex int
	sumGasUsed := uint64(0)
	if len(sorter) > 0 {
		sumGasUsed = sorter[0].gasUsed
	}
	for i, p := range rewardPercentiles {
		thresholdGasUsed := uint64(float64(blockGasUsed) * p / 100)
		for sumGasUsed < thresholdGasUsed && txIndex < tendermintTxCount-1 {
			txIndex++
			sumGasUsed += sorter[txIndex].gasUsed
		}

		chosenReward := big.NewInt(0)
		if 0 <= txIndex && txIndex < len(sorter) {
			chosenReward = sorter[txIndex].reward
		}
		targetOneFeeHistory.Reward[i] = chosenReward
	}

	return nil
}

type OneFeeHistory struct {
	BaseFee *big.Int
	Reward  []*big.Int
	GasUsed float64
}

func (e *EVMBackend) FeeHistory(userBlockCount rpc.DecimalOrHex, lastBlock rpc.BlockNumber, rewardPercentiles []float64) (*rpctypes.FeeHistoryResult, error) {

	var blockEnd int64 = int64(lastBlock)

	if blockEnd <= 0 {
		blockNumber, err := e.BlockNumber()
		if err != nil {
			return nil, err
		}
		blockEnd = int64(blockNumber)
	}
	var blockStart int64 = blockEnd - int64(userBlockCount)
	if blockStart < 0 {
		blockStart = 0
	}

	var blockCount int64 = blockEnd - blockStart

	var OldestBlock *hexutil.Big = (*hexutil.Big)(big.NewInt(blockStart))

	// prepare space
	var Reward [][]*hexutil.Big = make([][]*hexutil.Big, blockCount)
	var rewardcount = len(rewardPercentiles)
	for i := 0; i < int(blockCount); i++ {
		Reward[i] = make([]*hexutil.Big, rewardcount)
	}
	var BaseFee []*hexutil.Big = make([]*hexutil.Big, blockCount)
	var GasUsedRatio []float64 = make([]float64, blockCount)

	// fetch block
	for blockID := blockStart; blockID < blockEnd; blockID++ {
		index := int32(blockID - blockStart)
		// eth block
		ethBlock, err := e.GetBlockByNumber(rpctypes.BlockNumber(blockID), true)
		if err != nil {
			return nil, err
		}

		// tendermint block
		tendermintblock, tenderminterr := e.GetTendermintBlockByNumber(rpctypes.BlockNumber(blockID))
		if tenderminterr != nil {
			return nil, err
		}

		// tendermint block result
		tendermintBlockResult, err := e.clientCtx.Client.BlockResults(e.ctx, &tendermintblock.Block.Height)
		if err != nil {
			e.logger.Debug("EthBlockFromTendermint block result not found", "height", tendermintblock.Block.Height, "error", err.Error())
			return nil, err
		}

		var onefeehistory OneFeeHistory = OneFeeHistory{}
		processerr := e.processBlock(tendermintblock, &ethBlock, rewardPercentiles, tendermintBlockResult, &onefeehistory)
		if processerr != nil {
			return nil, processerr
		}

		// copy
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
}
