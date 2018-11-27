package client

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/x/stake/client/cli"
	"github.com/cosmos/cosmos-sdk/x/stake/client/rest"
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
	// Get all delegations from a delegator
	rs.Mux.HandleFunc(
		"/stake/delegators/{delegatorAddr}/delegations",
		rest.DelegatorDelegationsHandlerFn(rs.CliCtx, rs.Cdc),
	).Methods("GET")

	// Get all unbonding delegations from a delegator
	rs.Mux.HandleFunc(
		"/stake/delegators/{delegatorAddr}/unbonding_delegations",
		rest.DelegatorUnbondingDelegationsHandlerFn(rs.CliCtx, rs.Cdc),
	).Methods("GET")

	// Get all redelegations from a delegator
	rs.Mux.HandleFunc(
		"/stake/delegators/{delegatorAddr}/redelegations",
		rest.DelegatorRedelegationsHandlerFn(rs.CliCtx, rs.Cdc),
	).Methods("GET")

	// Get all staking txs (i.e msgs) from a delegator
	rs.Mux.HandleFunc(
		"/stake/delegators/{delegatorAddr}/txs",
		rest.DelegatorTxsHandlerFn(rs.CliCtx, rs.Cdc),
	).Methods("GET")

	// Query all validators that a delegator is bonded to
	rs.Mux.HandleFunc(
		"/stake/delegators/{delegatorAddr}/validators",
		rest.DelegatorValidatorsHandlerFn(rs.CliCtx, rs.Cdc),
	).Methods("GET")

	// Query a validator that a delegator is bonded to
	rs.Mux.HandleFunc(
		"/stake/delegators/{delegatorAddr}/validators/{validatorAddr}",
		rest.DelegatorValidatorHandlerFn(rs.CliCtx, rs.Cdc),
	).Methods("GET")

	// Query a delegation between a delegator and a validator
	rs.Mux.HandleFunc(
		"/stake/delegators/{delegatorAddr}/delegations/{validatorAddr}",
		rest.DelegationHandlerFn(rs.CliCtx, rs.Cdc),
	).Methods("GET")

	// Query all unbonding delegations between a delegator and a validator
	rs.Mux.HandleFunc(
		"/stake/delegators/{delegatorAddr}/unbonding_delegations/{validatorAddr}",
		rest.UnbondingDelegationHandlerFn(rs.CliCtx, rs.Cdc),
	).Methods("GET")

	// Get all validators
	rs.Mux.HandleFunc(
		"/stake/validators",
		rest.ValidatorsHandlerFn(rs.CliCtx, rs.Cdc),
	).Methods("GET")

	// Get a single validator info
	rs.Mux.HandleFunc(
		"/stake/validators/{validatorAddr}",
		rest.ValidatorHandlerFn(rs.CliCtx, rs.Cdc),
	).Methods("GET")

	// Get all delegations to a validator
	rs.Mux.HandleFunc(
		"/stake/validators/{validatorAddr}/delegations",
		rest.ValidatorDelegationsHandlerFn(rs.CliCtx, rs.Cdc),
	).Methods("GET")

	// Get all unbonding delegations from a validator
	rs.Mux.HandleFunc(
		"/stake/validators/{validatorAddr}/unbonding_delegations",
		rest.ValidatorUnbondingDelegationsHandlerFn(rs.CliCtx, rs.Cdc),
	).Methods("GET")

	// Get all outgoing redelegations from a validator
	rs.Mux.HandleFunc(
		"/stake/validators/{validatorAddr}/redelegations",
		rest.ValidatorRedelegationsHandlerFn(rs.CliCtx, rs.Cdc),
	).Methods("GET")

	// Get the current state of the staking pool
	rs.Mux.HandleFunc(
		"/stake/pool",
		rest.PoolHandlerFn(rs.CliCtx, rs.Cdc),
	).Methods("GET")

	// Get the current staking parameter values
	rs.Mux.HandleFunc(
		"/stake/parameters",
		rest.ParamsHandlerFn(rs.CliCtx, rs.Cdc),
	).Methods("GET")

	// Make a delegation transaction
	rs.Mux.HandleFunc(
		"/stake/delegators/{delegatorAddr}/delegations",
		rest.DelegationsRequestHandlerFn(rs.Cdc, rs.KeyBase, rs.CliCtx),
	).Methods("POST")
}

// GetQueryCmd returns the cli query commands for this module
func (mc ModuleClient) GetQueryCmd() *cobra.Command {
	stakeQueryCmd := &cobra.Command{
		Use:   "stake",
		Short: "Querying commands for the staking module",
	}
	stakeQueryCmd.AddCommand(client.GetCommands(
		cli.GetCmdQueryDelegation(mc.storeKey, mc.cdc),
		cli.GetCmdQueryDelegations(mc.storeKey, mc.cdc),
		cli.GetCmdQueryUnbondingDelegation(mc.storeKey, mc.cdc),
		cli.GetCmdQueryUnbondingDelegations(mc.storeKey, mc.cdc),
		cli.GetCmdQueryRedelegation(mc.storeKey, mc.cdc),
		cli.GetCmdQueryRedelegations(mc.storeKey, mc.cdc),
		cli.GetCmdQueryValidator(mc.storeKey, mc.cdc),
		cli.GetCmdQueryValidators(mc.storeKey, mc.cdc),
		cli.GetCmdQueryValidatorDelegations(mc.storeKey, mc.cdc),
		cli.GetCmdQueryValidatorUnbondingDelegations(mc.storeKey, mc.cdc),
		cli.GetCmdQueryValidatorRedelegations(mc.storeKey, mc.cdc),
		cli.GetCmdQueryParams(mc.storeKey, mc.cdc),
		cli.GetCmdQueryPool(mc.storeKey, mc.cdc))...)

	return stakeQueryCmd

}

// GetTxCmd returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	stakeTxCmd := &cobra.Command{
		Use:   "stake",
		Short: "Staking transaction subcommands",
	}

	stakeTxCmd.AddCommand(client.PostCommands(
		cli.GetCmdCreateValidator(mc.cdc),
		cli.GetCmdEditValidator(mc.cdc),
		cli.GetCmdDelegate(mc.cdc),
		cli.GetCmdRedelegate(mc.storeKey, mc.cdc),
		cli.GetCmdUnbond(mc.storeKey, mc.cdc),
	)...)

	return stakeTxCmd
}
