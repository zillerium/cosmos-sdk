package ibc

import (
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgOpenConn:
			return handleMsgOpenConn(ctx, k, msg)
		case MsgUpdateConn:
			return handleMsgUpdateConn(ctx, k, msg)
		default:
			errMsg := "Unrecognized IBC Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgOpenConn(ctx sdk.Context, k Keeper, msg MsgOpenConn) sdk.Result {
	r := k.connRuntime(ctx, msg.ChainID)

	if r.connEstablished() {
		return ErrConnAlreadyEstablished(k.codespace).Result()
	}

	height := uint64(msg.ROT.Height())
	r.setCommitHeight(height)
	r.setCommit(height, msg.ROT)

	return sdk.Result{}
}

func handleMsgUpdateConn(ctx sdk.Context, k Keeper, msg MsgUpdateConn) sdk.Result {
	r := k.connRuntime(ctx, msg.SrcChain)

	if !r.connEstablished() {
		return ErrConnNotEstablished(k.codespace).Result()
	}

	lastheight := r.getCommitHeight()
	height := uint64(msg.Commit.Commit.Height())
	if height < lastheight {
		return ErrInvalidHeight(k.codespace).Result()
	}

	// TODO: add lc verificatioon
	/*
		lastcommit := r.getCommit(lastheight)

		cert := lite.NewDynamicCertifier(msg.SrcChain, commit.Validators, height)
		if err := cert.Update(msg.Commit); err != nil {
			return ErrUpdateCommitFailed(k.codespace, err).Result()
		}

		k.setCommit(ctx, msg.SrcChain, msg.Commit.Height(), msg.Commit)
	*/
	r.setCommit(height, msg.Commit)
	return sdk.Result{}
}
