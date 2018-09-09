package rest

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/tendermint/tendermint/lite"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/ibc"
)

const (
	ChainID = "chain-id"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *wire.Codec, kb keys.Keybase) {
	r.HandleFunc(fmt.Sprintf("/ibc/conn/{%s}/open", ChainID), connOpenHandlerFn(cdc, kb, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/ibc/conn/{%s}/update", ChainID), connUpdateHandlerFn(cdc, kb, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/ibc/conn/{%s}/close", ChainID), connCloseHandler(cdc, kb, cliCtx)).Methods("POST")
}

type connOpenReq struct {
	LocalAccountName string `json:"name"`
	Password         string `json:"password"`
	ChainID          string `json:"chain_id"`
	AccountNumber    int64  `json:"account_number"`
	Sequence         int64  `json:"sequence"`
	Gas              int64  `json:"gas"`
	GasAdjustment    string `json:"gas_adjustment"`

	SrcChain string          `json:"src_chain"`
	ROT      lite.FullCommit `json:"root_of_trust"`
}

func connOpenHandlerFn(cdc *wire.Codec, kb keys.Keybase, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m connOpenReq
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = cdc.UnmarshalJSON(body, &m)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		info, err := kb.Get(m.LocalAccountName)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		msg := client.BuildMsg(sdk.AccAddress(info.GetPubKey().Address()))
	}
}
