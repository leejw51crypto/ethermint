package main

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
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
	out, _ := wallet0.Derive(hdpath, true)
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

			/*
						From     common.Address  `json:"from"`
				To       *common.Address `json:"to"`
				Gas      *hexutil.Uint64 `json:"gas"`
				GasPrice *hexutil.Big    `json:"gasPrice"`
				Value    *hexutil.Big    `json:"value"`
				Nonce    *hexutil.Uint64 `json:"nonce"`
				// We accept "data" and "input" for backwards-compatibility reasons. "input" is the
				// newer name and should be preferred by clients.
				Data  *hexutil.Bytes `json:"data"`
				Input *hexutil.Bytes `json:"input"`
				// For non-legacy transactions
				AccessList *ethtypes.AccessList `json:"accessList,omitempty"`
				ChainID    *hexutil.Big         `json:"chainId,omitempty"`
			*/
			test1, _ := hexutil.Decode("0x1234")
			test2 := hexutil.Encode(test1)
			fmt.Printf("test %v %s\n", test1, test2)

			//fromaddr2, _ := hexutil.Decode("0x3492dEc151Aa6179e13F775eD249185478F3D8ad")
			fromaddr2, _ := hexutil.Decode("0x48B212A71EBBB202F7CFD1AACEE3A36FDE2FBC51")
			fromaddr := common.BytesToAddress(fromaddr2)
			toaddr2, _ := hexutil.Decode("0x14F7B5EFAE3B1ECE60CB6AC0ACC67367E2C2E0F8")
			toaddr := common.BytesToAddress(toaddr2)

			gas := uint64(20000000000)
			//	nonce := uint64(0)

			data2, _ := hexutil.Decode("0x")
			data := hexutil.Bytes(data2)

			input2, _ := hexutil.Decode("0x")
			input := hexutil.Bytes(input2)
			fmt.Printf("%v %v\n", data, input)
			sendarg := rpctypes.SendTxArgs{
				From:     fromaddr,
				To:       &toaddr,
				Gas:      (*hexutil.Uint64)(&gas),
				GasPrice: (*hexutil.Big)(big.NewInt(20000000)),
				Value:    (*hexutil.Big)(big.NewInt(2)),
				Data:     nil,
				Input:    nil,
				//Data:     &data,
				//Input:    &input,
			}
			fmt.Printf("sendarg= %+v\n", sendarg)
			txhash, err := SendTransactionEth(clientCtx, evmBackend, queryClient, sendarg)
			fmt.Printf("txhash= %v\n", txhash)
			return err
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	cmd.AddCommand(addCmd, signCmd, sendCmd)

	return cmd
}
