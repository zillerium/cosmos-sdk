package ibc

import (
	"encoding/json"

	"github.com/tendermint/tendermint/lite"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgOpenConn defines the message that is used for open a c
// that receives msg from another chain
type MsgOpenConn struct {
	ROT     lite.FullCommit
	ChainID string
	Signer  sdk.AccAddress
}

func (msg MsgOpenConn) Type() string {
	return "ibc"
}

func (msg MsgOpenConn) Name() string {
	return "open_conn"
}

func (msg MsgOpenConn) GetSignBytes() []byte {
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return bz
}

func (msg MsgOpenConn) ValidateBasic() sdk.Error {
	if msg.ROT.Height() < 0 {
		// XXX: Codespace will be removed
		return ErrInvalidHeight(111)
	}
	return nil
}

func (msg MsgOpenConn) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

type MsgUpdateConn struct {
	SrcChain string
	Commit   lite.FullCommit
	//PacketProof
	Signer sdk.AccAddress
}

func (msg MsgUpdateConn) Type() string {
	return "ibc"
}

func (msg MsgUpdateConn) Name() string {
	return "update_conn"
}

func (msg MsgUpdateConn) GetSignBytes() []byte {
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return bz
}

func (msg MsgUpdateConn) ValidateBasic() sdk.Error {
	if msg.Commit.Commit.Height() < 0 {
		// XXX: Codespace will be removed
		return ErrInvalidHeight(111)
	}
	return nil
}

func (msg MsgUpdateConn) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}
