package client

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/x/slashing/client/cli"
	"github.com/cosmos/cosmos-sdk/x/slashing/client/rest"
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"
)

// ModuleClient exports all client functionality from this module
type ModuleClient struct {
	storeKey string
	cdc      *amino.Codec
}

func NewModuleClient(storeKey string, cdc *amino.Codec) ModuleClient {
	return ModuleClient{storeKey, cdc}
}

// RegisterRoutes registers staking-related REST handlers to a router
func (mc ModuleClient) RegisterRoutes(rs lcd.RestServer) {
	rs.Mux.HandleFunc(
		"/slashing/validators/{validatorPubKey}/signing_info",
		rest.SigningInfoHandlerFn(rs.CliCtx, mc.storeKey, rs.Cdc),
	).Methods("GET")

	rs.Mux.HandleFunc(
		"/slashing/validators/{validatorAddr}/unjail",
		rest.UnjailRequestHandlerFn(rs.Cdc, rs.KeyBase, rs.CliCtx),
	).Methods("POST")
}

// GetQueryCmd returns the cli query commands for this module
func (mc ModuleClient) GetQueryCmd() *cobra.Command {
	// Group slashing queries under a subcommand
	slashingQueryCmd := &cobra.Command{
		Use:   "slashing",
		Short: "Querying commands for the slashing module",
	}

	slashingQueryCmd.AddCommand(client.GetCommands(
		cli.GetCmdQuerySigningInfo(mc.storeKey, mc.cdc))...)

	return slashingQueryCmd

}

// GetTxCmd returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	slashingTxCmd := &cobra.Command{
		Use:   "slashing",
		Short: "Slashing transactions subcommands",
	}

	slashingTxCmd.AddCommand(client.PostCommands(
		cli.GetCmdUnjail(mc.cdc),
	)...)

	return slashingTxCmd
}
