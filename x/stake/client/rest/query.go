package rest

import (
	"net/http"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/stake/tags"

	"github.com/gorilla/mux"
)

const storeName = "stake"

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {

	// Get all delegations from a delegator
	r.HandleFunc(
		"/stake/delegators/{delegatorAddr}/delegations",
		DelegatorDelegationsHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Get all unbonding delegations from a delegator
	r.HandleFunc(
		"/stake/delegators/{delegatorAddr}/unbonding_delegations",
		DelegatorUnbondingDelegationsHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Get all redelegations from a delegator
	r.HandleFunc(
		"/stake/delegators/{delegatorAddr}/redelegations",
		DelegatorRedelegationsHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Get all staking txs (i.e msgs) from a delegator
	r.HandleFunc(
		"/stake/delegators/{delegatorAddr}/txs",
		DelegatorTxsHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Query all validators that a delegator is bonded to
	r.HandleFunc(
		"/stake/delegators/{delegatorAddr}/validators",
		DelegatorValidatorsHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Query a validator that a delegator is bonded to
	r.HandleFunc(
		"/stake/delegators/{delegatorAddr}/validators/{validatorAddr}",
		DelegatorValidatorHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Query a delegation between a delegator and a validator
	r.HandleFunc(
		"/stake/delegators/{delegatorAddr}/delegations/{validatorAddr}",
		DelegationHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Query all unbonding delegations between a delegator and a validator
	r.HandleFunc(
		"/stake/delegators/{delegatorAddr}/unbonding_delegations/{validatorAddr}",
		UnbondingDelegationHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Get all validators
	r.HandleFunc(
		"/stake/validators",
		ValidatorsHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Get a single validator info
	r.HandleFunc(
		"/stake/validators/{validatorAddr}",
		ValidatorHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Get all delegations to a validator
	r.HandleFunc(
		"/stake/validators/{validatorAddr}/delegations",
		ValidatorDelegationsHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Get all unbonding delegations from a validator
	r.HandleFunc(
		"/stake/validators/{validatorAddr}/unbonding_delegations",
		ValidatorUnbondingDelegationsHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Get all outgoing redelegations from a validator
	r.HandleFunc(
		"/stake/validators/{validatorAddr}/redelegations",
		ValidatorRedelegationsHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Get the current state of the staking pool
	r.HandleFunc(
		"/stake/pool",
		PoolHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Get the current staking parameter values
	r.HandleFunc(
		"/stake/parameters",
		ParamsHandlerFn(cliCtx, cdc),
	).Methods("GET")

}

// DelegatorDelegationsHandlerFn handles the /stake/delegators/{delegatorAddr}/delegations route
// HTTP request handler to query a delegator delegations
func DelegatorDelegationsHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return queryDelegator(cliCtx, cdc, "custom/stake/delegatorDelegations")
}

// DelegatorUnbondingDelegationsHandlerFn handles the /stake/delegators/{delegatorAddr}/unbonding_delegations route
// HTTP request handler to query a delegator unbonding delegations
func DelegatorUnbondingDelegationsHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return queryDelegator(cliCtx, cdc, "custom/stake/delegatorUnbondingDelegations")
}

// DelegatorRedelegationsHandlerFn handles the /stake/delegators/{delegatorAddr}/redelegations route
// HTTP request handler to query a delegator redelegations
func DelegatorRedelegationsHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return queryDelegator(cliCtx, cdc, "custom/stake/delegatorRedelegations")
}

// DelegatorTxsHandlerFn handles the /stake/delegators/{delegatorAddr}/txs route
// HTTP request handler to query all staking txs (msgs) from a delegator
func DelegatorTxsHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var typesQuerySlice []string
		vars := mux.Vars(r)
		delegatorAddr := vars["delegatorAddr"]

		_, err := sdk.AccAddressFromBech32(delegatorAddr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		node, err := cliCtx.GetNode()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Get values from query

		typesQuery := r.URL.Query().Get("type")
		trimmedQuery := strings.TrimSpace(typesQuery)
		if len(trimmedQuery) != 0 {
			typesQuerySlice = strings.Split(trimmedQuery, " ")
		}

		noQuery := len(typesQuerySlice) == 0
		isBondTx := contains(typesQuerySlice, "bond")
		isUnbondTx := contains(typesQuerySlice, "unbond")
		isRedTx := contains(typesQuerySlice, "redelegate")
		var txs = []tx.Info{}
		var actions []string

		switch {
		case isBondTx:
			actions = append(actions, string(tags.ActionDelegate))
		case isUnbondTx:
			actions = append(actions, string(tags.ActionBeginUnbonding))
			actions = append(actions, string(tags.ActionCompleteUnbonding))
		case isRedTx:
			actions = append(actions, string(tags.ActionBeginRedelegation))
			actions = append(actions, string(tags.ActionCompleteRedelegation))
		case noQuery:
			actions = append(actions, string(tags.ActionDelegate))
			actions = append(actions, string(tags.ActionBeginUnbonding))
			actions = append(actions, string(tags.ActionCompleteUnbonding))
			actions = append(actions, string(tags.ActionBeginRedelegation))
			actions = append(actions, string(tags.ActionCompleteRedelegation))
		default:
			w.WriteHeader(http.StatusNoContent)
			return
		}

		for _, action := range actions {
			foundTxs, errQuery := queryTxs(node, cliCtx, cdc, action, delegatorAddr)
			if errQuery != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			}
			txs = append(txs, foundTxs...)
		}

		res, err := cdc.MarshalJSON(txs)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

// UnbondingDelegationHandlerFn handles the /stake/delegators/{delegatorAddr}/unbonding_delegations/{validatorAddr} route
// HTTP request handler to query an unbonding-delegation
func UnbondingDelegationHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return queryBonds(cliCtx, cdc, "custom/stake/unbondingDelegation")
}

// DelegationHandlerFn handles the /stake/delegators/{delegatorAddr}/delegations/{validatorAddr} route
// HTTP request handler to query a delegation
func DelegationHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return queryBonds(cliCtx, cdc, "custom/stake/delegation")
}

// DelegatorValidatorsHandlerFn handles the /stake/delegators/{delegatorAddr}/validators route
// HTTP request handler to query all delegator bonded validators
func DelegatorValidatorsHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return queryDelegator(cliCtx, cdc, "custom/stake/delegatorValidators")
}

// DelegatorValidatorHandlerFn handles the /stake/delegators/{delegatorAddr}/validators/{validatorAddr} route
// HTTP request handler to get information from a currently bonded validator
func DelegatorValidatorHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return queryBonds(cliCtx, cdc, "custom/stake/delegatorValidator")
}

// ValidatorsHandlerFn handles the /stake/validators function
// HTTP request handler to query list of validators
func ValidatorsHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := cliCtx.QueryWithData("custom/stake/validators", nil)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

// ValidatorHandlerFn handles the /stake/validators/{validatorAddr} route
// HTTP request handler to query the validator information from a given validator address
func ValidatorHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return queryValidator(cliCtx, cdc, "custom/stake/validator")
}

// ValidatorDelegationsHandlerFn handles the /stake/validators/{validatorAddr}/delegations route
// HTTP request handler to query all unbonding delegations from a validator
func ValidatorDelegationsHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return queryValidator(cliCtx, cdc, "custom/stake/validatorDelegations")
}

// ValidatorUnbondingDelegationsHandlerFn handles the /stake/validators/{validatorAddr}/unbonding_delegations route
// HTTP request handler to query all unbonding delegations from a validator
func ValidatorUnbondingDelegationsHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return queryValidator(cliCtx, cdc, "custom/stake/validatorUnbondingDelegations")
}

// ValidatorRedelegationsHandlerFn handles the /stake/validators/{validatorAddr}/redelegations route
// HTTP request handler to query all redelegations from a source validator
func ValidatorRedelegationsHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return queryValidator(cliCtx, cdc, "custom/stake/validatorRedelegations")
}

// PoolHandlerFn handles the /stake/pool route
// HTTP request handler to query the pool information
func PoolHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := cliCtx.QueryWithData("custom/stake/pool", nil)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

// ParamsHandlerFn handles the /stake/parameters route
// HTTP request handler to query the staking params values
func ParamsHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := cliCtx.QueryWithData("custom/stake/parameters", nil)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
