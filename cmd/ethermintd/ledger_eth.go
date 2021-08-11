package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/pkg/errors"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"

	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"

	ethtypes "github.com/ethereum/go-ethereum/core/types"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/tharsis/ethermint/ethereum/rpc/backend"
	rpctypes "github.com/tharsis/ethermint/ethereum/rpc/types"
	ethermint "github.com/tharsis/ethermint/types"
	"github.com/tharsis/ethermint/usbwallet"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"
)

func myassert(check bool) {
	if !check {
		panic("")
	}
}

func mypadding(src []byte, dst []byte) {
	n := len(dst)

	for i := 0; i < n; i++ {
		dst[i] = 0
	}

	fmt.Printf("before src %v dst %v\n", src, dst)
	m := len(src)
	for i := 0; i < m; i++ {
		dst[n-1-i] = src[m-1-i]
	}

	fmt.Printf("after src %v dst %v\n", src, dst)
}

func SignMsg(msg *evmtypes.MsgEthereumTx, ethSigner ethtypes.Signer, keyringSigner keyring.Signer) error {

	fmt.Printf("this is my engine\n")
	// Start a USB hub for Ledger hardware wallets
	ledgerhub, err := usbwallet.NewLedgerHub()
	if err != nil {
		return nil
	}
	fmt.Printf("found ledger hub %v\n", ledgerhub)
	w := ledgerhub.Wallets()
	openerr := w[0].Open("")
	fmt.Printf("open %v\n", openerr)
	index := uint32(0)
	hdpath := []uint32{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0, index}

	out, _ := w[0].Derive(hdpath, true)
	fmt.Printf("----------------------------------------------\n")
	fmt.Printf("derived index %d = pubkey %d %s address %v\n", index, len(out.Pubkey), hexutil.Encode(out.Pubkey), out)

	wallet0 := w[0]
	accounts := wallet0.Accounts()
	fmt.Printf("accounts  length %d\n", len(accounts))
	for index, element := range accounts {
		fmt.Printf("index %d account %v\n", index, element)
	}

	// tx ------------------------------------------
	tx := msg.AsTransaction()
	txHash := ethSigner.Hash(tx)

	// sign ---------------------------------------------------
	key, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(key.PublicKey)
	fmt.Printf("addr %v\n", addr)
	//chainid := big.NewInt(2)
	chainid := ethSigner.ChainID()
	//txfound, err := wallet0.SignTx(accounts[0], mytypes.NewTransaction(0, addr, new(big.Int), 0, new(big.Int), nil), chainid)
	txfound, err := wallet0.SignTx(accounts[0], tx, chainid)
	txjson, _ := txfound.MarshalJSON()
	fmt.Printf("tx json %v\n", string(txjson))
	v2, r2, s2 := txfound.RawSignatureValues()
	v := hexutil.EncodeBig(v2)
	r := hexutil.EncodeBig(r2)
	s := hexutil.EncodeBig(s2)

	v3, _ := hexutil.Decode(v)

	// pading
	r3a, _ := hexutil.Decode(r)
	s3a, _ := hexutil.Decode(s)

	// test code
	r3 := make([]byte, 32, 32)
	s3 := make([]byte, 32, 32)
	mypadding(r3a, r3)
	mypadding(s3a, s3)

	myassert(len(v3) == 1)
	myassert(len(r3) == 32)
	myassert(len(s3) == 32)

	hwsignature := append(r3, s3...)
	hwsignature = append(hwsignature, v3...)
	myassert(len(hwsignature) == 65)

	fmt.Printf("## ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")
	fmt.Printf("signature %d %+v\n", len(hwsignature), hwsignature)
	// order  r(32), s(32), v(1)
	/*
			signature = signature.substr(2); //remove 0x
		const r = '0x' + signature.slice(0, 64)
		const s = '0x' + signature.slice(64, 128)
		const v = '0x' + signature.slice(128, 130)
	*/

	fmt.Printf("####  SignMsg ~~~~~~~~~~~~~~~~\n")
	from := msg.GetFrom()
	if from.Empty() {
		return fmt.Errorf("sender address not defined for message")
	}

	sig, _, err := keyringSigner.SignByAddress(from, txHash.Bytes())
	fmt.Printf("## ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")
	fmt.Printf("signature %d %+v\n", len(sig), sig)
	if err != nil {
		return err
	}

	tx, err = tx.WithSignature(ethSigner, sig)
	if err != nil {
		return err
	}

	msg.FromEthereumTx(tx)
	return nil
}

// --------------------------------------------------------------------------

func getAccountNonce(
	clientCtx client.Context,
	backend backend.Backend,
	chainIDEpoch *big.Int,
	accAddr common.Address, pending bool, height int64) (uint64, error) {
	queryClient := authtypes.NewQueryClient(clientCtx)
	res, err := queryClient.Account(rpctypes.ContextWithHeight(height), &authtypes.QueryAccountRequest{Address: sdk.AccAddress(accAddr.Bytes()).String()})
	if err != nil {
		return 0, err
	}
	var acc authtypes.AccountI
	if err := clientCtx.InterfaceRegistry.UnpackAny(res.Account, &acc); err != nil {
		return 0, err
	}

	nonce := acc.GetSequence()

	if !pending {
		return nonce, nil
	}

	// the account retriever doesn't include the uncommitted transactions on the nonce so we need to
	// to manually add them.
	pendingTxs, err := backend.PendingTransactions()
	if err != nil {
		return nonce, nil
	}

	// add the uncommitted txs to the nonce counter
	// only supports `MsgEthereumTx` style tx
	for _, tx := range pendingTxs {
		msg, err := evmtypes.UnwrapEthereumMsg(tx)
		if err != nil {
			// not ethereum tx
			continue
		}

		sender, err := msg.GetSender(chainIDEpoch)
		if err != nil {
			continue
		}
		if sender == accAddr {
			nonce++
		}
	}

	return nonce, nil
}

// EstimateGas returns an estimate of gas usage for the given smart contract call.
func EstimateGas(clientCtx client.Context, args evmtypes.CallArgs, blockNrOptional *rpctypes.BlockNumber) (hexutil.Uint64, error) {
	queryClient := rpctypes.NewQueryClient(clientCtx)
	blockNr := rpctypes.EthPendingBlockNumber
	if blockNrOptional != nil {
		blockNr = *blockNrOptional
	}

	bz, err := json.Marshal(&args)
	if err != nil {
		return 0, err
	}
	req := evmtypes.EthCallRequest{Args: bz, GasCap: ethermint.DefaultRPCGasLimit}

	// From ContextWithHeight: if the provided height is 0,
	// it will return an empty context and the gRPC query will use
	// the latest block height for querying.
	res, err := queryClient.EstimateGas(rpctypes.ContextWithHeight(blockNr.Int64()), &req)
	if err != nil {
		return 0, err
	}
	return hexutil.Uint64(res.Gas), nil
}

func setTxDefaults(clientCtx client.Context,
	backend backend.Backend,
	chainIDEpoch *big.Int,
	args rpctypes.SendTxArgs) (rpctypes.SendTxArgs, error) {

	if args.GasPrice == nil {
		// TODO: Change to either:
		// - min gas price from context once available through server/daemon, or
		// - suggest a gas price based on the previous included txs
		args.GasPrice = (*hexutil.Big)(big.NewInt(ethermint.DefaultGasPrice))
	}

	if args.Nonce == nil {
		// get the nonce from the account retriever
		// ignore error in case tge account doesn't exist yet
		nonce, _ := getAccountNonce(clientCtx, backend, chainIDEpoch, args.From, true, 0)
		//nonce := uint64(0)
		args.Nonce = (*hexutil.Uint64)(&nonce)
	}

	if args.Data != nil && args.Input != nil && !bytes.Equal(*args.Data, *args.Input) {
		return args, errors.New("both 'data' and 'input' are set and not equal. Please use 'input' to pass transaction call data")
	}

	if args.To == nil {
		// Contract creation
		var input []byte
		if args.Data != nil {
			input = *args.Data
		} else if args.Input != nil {
			input = *args.Input
		}

		if len(input) == 0 {
			return args, errors.New(`contract creation without any data provided`)
		}
	}

	if args.Gas == nil {
		// For backwards-compatibility reason, we try both input and data
		// but input is preferred.
		input := args.Input
		if input == nil {
			input = args.Data
		}

		callArgs := evmtypes.CallArgs{
			From:       &args.From, // From shouldn't be nil
			To:         args.To,
			Gas:        args.Gas,
			GasPrice:   args.GasPrice,
			Value:      args.Value,
			Data:       input,
			AccessList: args.AccessList,
		}
		blockNr := rpctypes.NewBlockNumber(big.NewInt(0))
		estimated, err := EstimateGas(clientCtx, callArgs, &blockNr)

		if err != nil {
			return args, err
		}
		//estimated := hexutil.Uint64(0)
		args.Gas = &estimated

	}

	if args.ChainID == nil {
		args.ChainID = (*hexutil.Big)(chainIDEpoch)
	}

	return args, nil
}

func SendTransactionEth(
	clientCtx client.Context,
	backend backend.Backend,
	queryClient *rpctypes.QueryClient,
	args rpctypes.SendTxArgs) (common.Hash, error) {

	epoch, epocherr := ethermint.ParseChainID(clientCtx.ChainID)
	fmt.Printf("epoch = %+v\n", epoch)
	if epocherr != nil {
		panic(epocherr)
	}

	// Look up the wallet containing the requested signer
	keyringinfo, err := clientCtx.Keyring.KeyByAddress(sdk.AccAddress(args.From.Bytes()))
	if err != nil {

		return common.Hash{}, fmt.Errorf("%s; %s", keystore.ErrNoMatch, err.Error())
	}
	fmt.Printf("keyring info %+v\n", keyringinfo)

	args, err = setTxDefaults(clientCtx, backend, epoch, args)
	if err != nil {
		return common.Hash{}, err
	}

	msg := args.ToTransaction()

	if err := msg.ValidateBasic(); err != nil {

		return common.Hash{}, err
	}

	// TODO: get from chain config
	signer := ethtypes.LatestSignerForChainID(args.ChainID.ToInt())

	fmt.Printf("################################\n")

	fmt.Printf("message to sign= %+v\n", msg.AsTransaction())

	// Sign transaction
	if err := SignMsg(msg, signer, clientCtx.Keyring); err != nil {

		return common.Hash{}, err
	}

	// Assemble transaction from fields
	builder, ok := clientCtx.TxConfig.NewTxBuilder().(authtx.ExtensionOptionsTxBuilder)
	if !ok {
		return common.Hash{}, err
	}

	option, err := codectypes.NewAnyWithValue(&evmtypes.ExtensionOptionsEthereumTx{})
	if err != nil {
		return common.Hash{}, err
	}

	builder.SetExtensionOptions(option)
	err = builder.SetMsgs(msg)
	if err != nil {
		return common.Hash{}, err
	}

	// Query params to use the EVM denomination
	res, err := queryClient.QueryClient.Params(context.Background(), &evmtypes.QueryParamsRequest{})
	if err != nil {
		return common.Hash{}, err
	}

	txData, err := evmtypes.UnpackTxData(msg.Data)
	if err != nil {
		return common.Hash{}, err
	}

	fees := sdk.Coins{sdk.NewCoin(res.Params.EvmDenom, sdk.NewIntFromBigInt(txData.Fee()))}
	builder.SetFeeAmount(fees)
	builder.SetGasLimit(msg.GetGas())

	// Encode transaction by default Tx encoder
	txEncoder := clientCtx.TxConfig.TxEncoder()
	txBytes, err := txEncoder(builder.GetTx())
	if err != nil {
		return common.Hash{}, err
	}

	txHash := msg.AsTransaction().Hash()

	// Broadcast transaction in sync mode (default)
	// NOTE: If error is encountered on the node, the broadcast will not return an error
	syncCtx := clientCtx.WithBroadcastMode(flags.BroadcastSync)
	rsp, err := syncCtx.BroadcastTx(txBytes)
	if err != nil || rsp.Code != 0 {
		if err == nil {
			err = errors.New(rsp.RawLog)
		}
		return txHash, err
	}

	// Return transaction hash
	return txHash, nil
}
