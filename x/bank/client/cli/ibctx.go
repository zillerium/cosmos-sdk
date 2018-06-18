package cli

import (
	"encoding/hex"
	"os"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	codec "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	"github.com/cosmos/cosmos-sdk/x/ibc"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/x/bank"
)

// IBC transfer command
func IBCSendTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "transfer",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithCodec(cdc)
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			msg, err := buildMsg(from)
			if err != nil {
				return err
			}
			if cliCtx.GenerateOnly {
				return utils.PrintUnsignedStdTx(txBldr, cliCtx, []sdk.Msg{msg})
			}

			return utils.SendTx(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(FlagTo, "", "Address to send coins")
	cmd.Flags().String(FlagAmount, "", "Amount of coins to send")
	cmd.Flags().String(FlagDestChain, "", "Destination chain to send coins")
	return cmd
}

func buildMsg(from sdk.AccAddress) (sdk.Msg, error) {
	amount := viper.GetString(FlagAmount)
	coins, err := sdk.ParseCoins(amount)
	if err != nil {
		return nil, err
	}

	dest := viper.GetString(FlagTo)
	bz, err := hex.DecodeString(dest)
	if err != nil {
		return nil, err
	}
	to := sdk.AccAddress(bz)

	payload := bank.PayloadSend{
		SrcAddr:  from,
		DestAddr: to,
		Coins:    coins,
	}

	msg := bank.MsgIBCSend{
		PayloadSend: payload,
		DestChain:   viper.GetString(FlagDestChain),
	}

	return msg, nil
}
