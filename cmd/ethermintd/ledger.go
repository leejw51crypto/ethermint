package main

import (
	"fmt"
	"os"

	"github.com/tharsis/ethermint/usbwallet"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clienttx "github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/spf13/cobra"
)

const (
	flagIndex = "index"
)

func runAddCmd(cmd *cobra.Command, _ []string) error {
	index, _ := cmd.Flags().GetUint32(flagIndex)
	fmt.Printf("add address index= %d\n", index)

	fmt.Printf("add ledger wallet index %d\n", index)
	ledgerhub, detecterr := usbwallet.NewLedgerHub()
	if detecterr != nil {
		fmt.Printf("ledger detect error %v\n", detecterr)
		return detecterr
	}
	w := ledgerhub.Wallets()
	wallet0 := w[0]
	openerr := wallet0.Open("")
	if openerr != nil {
		fmt.Printf("ledger open error %v\n", openerr)
		return openerr
	}

	// bip44, coin type, account, change ,index
	hdpath := []uint32{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0, index}
	out, _ := w[0].Derive(hdpath, true)
	fmt.Printf("Ledger Address Index= %d   Address= %s\n", index, out.Address.Hex())
	closeerr := wallet0.Close()
	if closeerr != nil {
		fmt.Printf("ledger close error %v\n", closeerr)
		return openerr
	}

	return nil
}

func ledgerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "ledger",
		Aliases:                    []string{"l"},
		Short:                      "Ledger subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	// add
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Add ledger address",
		RunE:  runAddCmd,
	}
	addCmd.Flags().Uint32(flagIndex, 0, "Address index number for HD derivation")

	// sign
	signCmd := &cobra.Command{
		Use:   "sign",
		Short: "Sign via ledger",
		RunE: func(_ *cobra.Command, _ []string) error {
			fmt.Println("sign transaction via ledger")
			return nil
		},
	}

	// send
	sendCmd := &cobra.Command{
		Use: "send [from_key_or_address] [to_address] [amount]",
		Short: `Send funds from one account to another. Note, the'--from' flag is
ignored as it is implied from [from_key_or_address].`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Flags().Set(flags.FlagFrom, args[0])
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			toAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			coins, err := sdk.ParseCoinsNormalized(args[2])
			if err != nil {
				return err
			}

			msg := types.NewMsgSend(clientCtx.GetFromAddress(), toAddr, coins)

			//return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
			txf := clienttx.NewFactoryCLI(clientCtx, cmd.Flags())

			return LedgerBroadcastTx(clientCtx, txf, msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	cmd.AddCommand(addCmd, signCmd, sendCmd)

	return cmd
}

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

func Sign(txconfig client.TxConfig, txf clienttx.Factory, name string, txBuilder client.TxBuilder, overwriteSig bool) error {
	/*if txf.keybase == nil {
		return errors.New("keybase must be set prior to signing a transaction")
	}*/

	signMode := txf.SignMode()
	if signMode == signing.SignMode_SIGN_MODE_UNSPECIFIED {
		// use the SignModeHandler's default mode if unspecified
		signMode = txconfig.SignModeHandler().DefaultMode()
	}

	/*if err := checkMultipleSigners(signMode, txBuilder.GetTx()); err != nil {
		return err
	}*/
	fmt.Printf("sign mode %v\n", signMode)
	key, err := txf.Keybase().Key(name)
	if err != nil {
		return err
	}
	pubKey := key.GetPubKey()

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

	// Generate the bytes to be signed.
	bytesToSign, err := txconfig.SignModeHandler().GetSignBytes(signMode, signerData, txBuilder.GetTx())
	if err != nil {
		return err
	}

	// Sign those bytes

	sigBytes, _, err := txf.Keybase().Sign(name, bytesToSign)
	if err != nil {
		return err
	}

	// Construct the SignatureV2 struct
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
