package main

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	"github.com/tharsis/ethermint/version"
)

func ledgerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "ledger",
		Aliases:                    []string{"l"},
		Short:                      "Ledger subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Print version info",
		RunE: func(_ *cobra.Command, _ []string) error {
			fmt.Println(version.Version())
			return nil
		},
	}

	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Add ledger address",
		RunE: func(_ *cobra.Command, _ []string) error {
			fmt.Println("add address")
			return nil
		},
	}

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

	cmd.AddCommand(infoCmd, addCmd, signCmd, sendCmd)

	return cmd
}
