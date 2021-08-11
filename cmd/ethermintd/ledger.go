package main

import (
	"fmt"

	"github.com/tharsis/ethermint/ethereum/rpc/backend"
	"github.com/tharsis/ethermint/usbwallet"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/log"
	rpctypes "github.com/tharsis/ethermint/ethereum/rpc/types"
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
			fmt.Printf("msg =%v\n", msg)

			evmBackend := backend.NewEVMBackend(log.NewNopLogger(), clientCtx)
			queryClient := rpctypes.NewQueryClient(clientCtx)

			sendarg := rpctypes.SendTxArgs{}
			txhash, err := SendTransactionEth(clientCtx, evmBackend, queryClient, sendarg)
			fmt.Printf("txhash= %v\n", txhash)
			return err
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	cmd.AddCommand(addCmd, signCmd, sendCmd)

	return cmd
}
