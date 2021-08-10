package main

import (
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tharsis/ethermint/usbwallet"

	"github.com/cosmos/cosmos-sdk/client"
	clienttx "github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

func LedgerBroadcastTx(clientCtx client.Context, txf clienttx.Factory, msgs ...sdk.Msg) error {

	txf, err := prepareFactory(clientCtx, txf)
	if err != nil {
		return err
	}

	if txf.SimulateAndExecute() || clientCtx.Simulate {
		_, adjusted, err := clienttx.CalculateGas(clientCtx, txf, msgs...)
		if err != nil {
			return err
		}

		txf = txf.WithGas(adjusted)
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", clienttx.GasEstimateResponse{GasEstimate: txf.Gas()})
	}

	if clientCtx.Simulate {
		return nil
	}

	tx, err := clienttx.BuildUnsignedTx(txf, msgs...)
	if err != nil {
		return err
	}

	/*
		if !clientCtx.SkipConfirm {
			out, err := clientCtx.TxConfig.TxJSONEncoder()(tx.GetTx())
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintf(os.Stderr, "%s\n\n", out)

			buf := bufio.NewReader(os.Stdin)
			ok, err := input.GetConfirmation("confirm transaction before signing and broadcasting", buf, os.Stderr)

			if err != nil || !ok {
				_, _ = fmt.Fprintf(os.Stderr, "%s\n", "cancelled transaction")
				return err
			}
		}
	*/

	tx.SetFeeGranter(clientCtx.GetFeeGranterAddress())
	//	err = clienttx.Sign(txf, clientCtx.GetFromName(), tx, true)

	err = Sign(clientCtx.TxConfig, txf, clientCtx.GetFromName(), tx, true)

	if err != nil {
		return err
	}

	txBytes, err := clientCtx.TxConfig.TxEncoder()(tx.GetTx())
	if err != nil {
		return err
	}

	// broadcast to a Tendermint node
	res, err := clientCtx.BroadcastTx(txBytes)
	if err != nil {
		return err
	}

	return clientCtx.PrintProto(res)
}

// prepareFactory ensures the account defined by ctx.GetFromAddress() exists and
// if the account number and/or the account sequence number are zero (not set),
// they will be queried for and set on the provided Factory. A new Factory with
// the updated fields will be returned.
func prepareFactory(clientCtx client.Context, txf clienttx.Factory) (clienttx.Factory, error) {
	from := clientCtx.GetFromAddress()

	if err := txf.AccountRetriever().EnsureExists(clientCtx, from); err != nil {
		return txf, err
	}

	initNum, initSeq := txf.AccountNumber(), txf.Sequence()
	if initNum == 0 || initSeq == 0 {
		num, seq, err := txf.AccountRetriever().GetAccountNumberSequence(clientCtx, from)
		if err != nil {
			return txf, err
		}

		if initNum == 0 {
			txf = txf.WithAccountNumber(num)
		}

		if initSeq == 0 {
			txf = txf.WithSequence(seq)
		}
	}

	return txf, nil
}

func checkMultipleSigners(mode signing.SignMode, tx authsigning.Tx) error {
	if mode == signing.SignMode_SIGN_MODE_DIRECT &&
		len(tx.GetSigners()) > 1 {
		return fmt.Errorf("Signing in DIRECT mode is only supported for transactions with one signer only")
	}
	return nil
}

func Sign(txconfig client.TxConfig, txf clienttx.Factory, name string, txBuilder client.TxBuilder, overwriteSig bool) error {
	/*if txf.keybase == nil {
		return errors.New("keybase must be set prior to signing a transaction")
	}*/

	signMode := txf.SignMode()
	if signMode == signing.SignMode_SIGN_MODE_UNSPECIFIED {
		// use the SignModeHandler's default mode if unspecified
		signMode = txconfig.SignModeHandler().DefaultMode()
	}

	if err := checkMultipleSigners(signMode, txBuilder.GetTx()); err != nil {
		return err
	}

	fmt.Printf("sign mode %v\n", signMode)
	key, err := txf.Keybase().Key(name)
	if err != nil {
		return err
	}
	pubKey := key.GetPubKey()
	// 33 bytes
	fmt.Printf("pubkey length %d  %s\n", len(pubKey.Bytes()), hexutil.Encode(pubKey.Bytes()))

	//	pubKey := make([]byte, 64, 64)

	signerData := authsigning.SignerData{
		ChainID:       txf.ChainID(),
		AccountNumber: txf.AccountNumber(),
		Sequence:      txf.Sequence(),
	}

	// For SIGN_MODE_DIRECT, calling SetSignatures calls setSignerInfos on
	// TxBuilder under the hood, and SignerInfos is needed to generated the
	// sign bytes. This is the reason for setting SetSignatures here, with a
	// nil signature.
	//
	// Note: this line is not needed for SIGN_MODE_LEGACY_AMINO, but putting it
	// also doesn't affect its generated sign bytes, so for code's simplicity
	// sake, we put it here.
	sigData := signing.SingleSignatureData{
		SignMode:  signMode,
		Signature: nil,
	}
	sig := signing.SignatureV2{
		PubKey:   pubKey,
		Data:     &sigData,
		Sequence: txf.Sequence(),
	}
	var prevSignatures []signing.SignatureV2
	if !overwriteSig {
		prevSignatures, err = txBuilder.GetTx().GetSignaturesV2()
		if err != nil {
			return err
		}
	}
	if err := txBuilder.SetSignatures(sig); err != nil {
		return err
	}

	bytesToSign, err := txconfig.SignModeHandler().GetSignBytes(signMode, signerData, txBuilder.GetTx())
	if err != nil {
		return err
	}

	fmt.Printf("bytes to sign length %d  %v\n", len(bytesToSign), hexutil.Encode(bytesToSign))
	digestBz := ethcrypto.Keccak256Hash(bytesToSign).Bytes()
	fmt.Printf("digest bytes %d  %v\n", len(digestBz), hexutil.Encode(digestBz))
	sigBytes, _, err := txf.Keybase().Sign(name, bytesToSign)
	ledgerSign(bytesToSign)

	// 65 bytes
	fmt.Printf("signature length %d  %s\n", len(sigBytes), hexutil.Encode(sigBytes))
	if err != nil {
		return err
	}

	sigData = signing.SingleSignatureData{
		SignMode:  signMode,
		Signature: sigBytes,
	}
	sig = signing.SignatureV2{
		PubKey:   pubKey,
		Data:     &sigData,
		Sequence: txf.Sequence(),
	}

	if overwriteSig {
		return txBuilder.SetSignatures(sig)
	}
	prevSignatures = append(prevSignatures, sig)
	return txBuilder.SetSignatures(prevSignatures...)
}

func ledgerSign(bytesToSign []byte) error {
	fmt.Printf("########################### ledger sign\n")
	fmt.Printf("ledger sign  %d   %s\n", len(bytesToSign), hexutil.Encode(bytesToSign))

	fmt.Printf("this is my engine\n")
	// Start a USB hub for Ledger hardware wallets
	ledgerhub, err := usbwallet.NewLedgerHub()
	if err != nil {
		return fmt.Errorf("Failed to start Ledger hub, disabling: %v", err)
	}
	fmt.Printf("found ledger hub %v\n", ledgerhub)
	w := ledgerhub.Wallets()
	fmt.Printf("wallets %+v\n", w)

	fmt.Printf("wallets length %d\n", len(w))
	openerr := w[0].Open("")
	fmt.Printf("open %v\n", openerr)
	index := uint32(0)
	// 44, coin   ,   account, change, index
	hdpath := []uint32{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0, index}

	out, _ := w[0].Derive(hdpath, true)
	fmt.Printf("derived index %d = %v,  %v\n", index, out.Address, out)

	wallet0 := w[0]
	accounts := wallet0.Accounts()
	fmt.Printf("accounts  length %d\n", len(accounts))

	for index, element := range accounts {
		fmt.Printf("index %d account %v\n", index, element)
	}

	key, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(key.PublicKey)
	fmt.Printf("address %v\n", addr)

	chainid := big.NewInt(18)
	tx, err := wallet0.SignTx(accounts[0], ethtypes.NewTransaction(0, addr, new(big.Int), 0, new(big.Int), nil), chainid)
	txjson, _ := tx.MarshalJSON()
	fmt.Printf("tx json %v\n", string(txjson))
	v, r, s := tx.RawSignatureValues()
	fmt.Printf("tx chainid %v %v %v\n", chainid, tx, err)
	fmt.Printf("signature v=%v r=%v s=%v\n", v, r, s)

	return nil
}
