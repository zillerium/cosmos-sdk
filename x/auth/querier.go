package auth

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// query endpoints supported by the auth Querier
const (
	QueryAccount = "account"
)

// creates a querier for auth REST endpoints
func NewQuerier(keeper AccountKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		ptr, querier := pathToQuerier(path, keeper)
		if querier == nil {
			return nil, sdk.ErrUnknownRequest("unknown auth query endpoint")
		}
		rawerr := keeper.cdc.UnmarshalJSON(req.Data, ptr)
		if rawerr := keeper.cdc.UnmarshalJSON(req.Data, ptr); rawerr != nil {
			return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", rawerr))
		}
		res, err := querier(ctx)
		if err != nil {
			return nil, err
		}
		bz, rawerr := codec.MarshalJSONIndent(keeper.cdc, res)
		if rawerr != nil {
			return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", rawerr.Error()))
		}

		return bz, nil

	}
}

type querier func(ctx sdk.Context) (interface{}, sdk.Error)

func pathToQuerier(path []string, k AccountKeeper) (sdk.QueryInput, querier) {
	switch path[0] {
	case QueryAccount:
		ptr := new(QueryAccountParams)
		return ptr, querierAccount(ptr, k)
	default:
		return nil, nil
	}
}

// defines the params for query: "custom/acc/account"
type QueryAccountParams struct {
	Address sdk.AccAddress
}

func (params QueryAccountParams) ValidateInput() sdk.Error {
	if params.Address.Empty() {
		return sdk.ErrInvalidAddress("missing query address")
	}
	return nil
}

func NewQueryAccountParams(addr sdk.AccAddress) QueryAccountParams {
	return QueryAccountParams{
		Address: addr,
	}
}

func querierAccount(params *QueryAccountParams, k AccountKeeper) querier {
	return func(ctx sdk.Context) (interface{}, sdk.Error) {
		account := k.GetAccount(ctx, params.Address)
		if account == nil {
			return nil, sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", params.Address))
		}
		return account, nil
	}
}
