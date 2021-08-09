package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/usbwallet"

	"github.com/cosmos/cosmos-sdk/client"
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

	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Add ledger address",
		RunE:  runAddCmd,
	}
	addCmd.Flags().Uint32(flagIndex, 0, "Address index number for HD derivation")

	signCmd := &cobra.Command{
		Use:   "sign",
		Short: "Sign via ledger",
		RunE: func(_ *cobra.Command, _ []string) error {
			fmt.Println("sign transaction via ledger")
			return nil
		},
	}

	sendCmd := &cobra.Command{
		Use:   "send",
		Short: "Send via ledger",
		RunE: func(_ *cobra.Command, _ []string) error {
			fmt.Println("add address via ledger")
			return nil
		},
	}

	cmd.AddCommand(addCmd, signCmd, sendCmd)

	return cmd
}
