package client

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	cli "github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	rest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
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

// RegisterRoutes exports the lite client route functionality that is exposed by this module
func (mc ModuleClient) RegisterRoutes(rs lcd.RestServer) {
	// Make a proposal
	rs.Mux.HandleFunc(
		"/gov/proposals",
		rest.PostProposalHandlerFn(rs.Cdc, rs.CliCtx),
	).Methods("POST")

	// Make a deposit for a proposal
	rs.Mux.HandleFunc(
		fmt.Sprintf("/gov/proposals/{%s}/deposits", rest.RestProposalID),
		rest.DepositHandlerFn(rs.Cdc, rs.CliCtx),
	).Methods("POST")

	// Vote on a proposal
	rs.Mux.HandleFunc(
		fmt.Sprintf("/gov/proposals/{%s}/votes", rest.RestProposalID),
		rest.VoteHandlerFn(rs.Cdc, rs.CliCtx),
	).Methods("POST")

	// Get the parameters for governance
	rs.Mux.HandleFunc(
		fmt.Sprintf("/gov/parameters/{%s}", rest.RestParamsType),
		rest.QueryParamsHandlerFn(rs.Cdc, rs.CliCtx),
	).Methods("GET")

	// Get a list of gov proposals
	rs.Mux.HandleFunc(
		"/gov/proposals",
		rest.QueryProposalsWithParameterFn(rs.Cdc, rs.CliCtx),
	).Methods("GET")

	// Get details for a single proposal
	rs.Mux.HandleFunc(
		fmt.Sprintf("/gov/proposals/{%s}", rest.RestProposalID),
		rest.QueryProposalHandlerFn(rs.Cdc, rs.CliCtx),
	).Methods("GET")

	// Get the deposits for a single proposal
	rs.Mux.HandleFunc(
		fmt.Sprintf("/gov/proposals/{%s}/deposits", rest.RestProposalID),
		rest.QueryDepositsHandlerFn(rs.Cdc, rs.CliCtx),
	).Methods("GET")

	// Get the details for an individual deposit on a proposal
	rs.Mux.HandleFunc(
		fmt.Sprintf("/gov/proposals/{%s}/deposits/{%s}", rest.RestProposalID, rest.RestDepositor),
		rest.QueryDepositHandlerFn(rs.Cdc, rs.CliCtx),
	).Methods("GET")

	// Get the tally of votes on a proposal
	rs.Mux.HandleFunc(
		fmt.Sprintf("/gov/proposals/{%s}/tally", rest.RestProposalID),
		rest.QueryTallyOnProposalHandlerFn(rs.Cdc, rs.CliCtx),
	).Methods("GET")

	// Get the list of votes on a proposal
	rs.Mux.HandleFunc(
		fmt.Sprintf("/gov/proposals/{%s}/votes", rest.RestProposalID),
		rest.QueryVotesOnProposalHandlerFn(rs.Cdc, rs.CliCtx),
	).Methods("GET")

	// Get the details for an individual vote
	rs.Mux.HandleFunc(
		fmt.Sprintf("/gov/proposals/{%s}/votes/{%s}", rest.RestProposalID, rest.RestVoter),
		rest.QueryVoteHandlerFn(rs.Cdc, rs.CliCtx),
	).Methods("GET")
}

// GetQueryCmd returns the cli query commands for this module
func (mc ModuleClient) GetQueryCmd() *cobra.Command {
	// Group gov queries under a subcommand
	govQueryCmd := &cobra.Command{
		Use:   "gov",
		Short: "Querying commands for the governance module",
	}

	govQueryCmd.AddCommand(client.GetCommands(
		cli.GetCmdQueryProposal(mc.storeKey, mc.cdc),
		cli.GetCmdQueryProposals(mc.storeKey, mc.cdc),
		cli.GetCmdQueryVote(mc.storeKey, mc.cdc),
		cli.GetCmdQueryVotes(mc.storeKey, mc.cdc),
		cli.GetCmdQueryParams(mc.storeKey, mc.cdc),
		cli.GetCmdQueryDeposit(mc.storeKey, mc.cdc),
		cli.GetCmdQueryDeposits(mc.storeKey, mc.cdc),
		cli.GetCmdQueryTally(mc.storeKey, mc.cdc))...)

	return govQueryCmd
}

// GetTxCmd returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	govTxCmd := &cobra.Command{
		Use:   "gov",
		Short: "Governance transactions subcommands",
	}

	govTxCmd.AddCommand(client.PostCommands(
		cli.GetCmdDeposit(mc.cdc),
		cli.GetCmdVote(mc.cdc),
		cli.GetCmdSubmitProposal(mc.cdc),
	)...)

	return govTxCmd
}
