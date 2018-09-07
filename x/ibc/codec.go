package ibc

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSend{}, "cosmos-sdk/ibc/Send", nil)
	cdc.RegisterConcrete(MsgReceive{}, "cosmos-sdk/ibc/Receive", nil)
	cdc.RegisterConcrete(MsgCleanup{}, "cosmos-sdk/ibc/Cleanup", nil)
	cdc.RegisterConcrete(MsgOpenConn{}, "cosmos-sdk/ibc/OpenConn", nil)
	cdc.RegisterConcrete(MsgUpdateConn{}, "cosmos-sdk/ibc/UpdateConn", nil)

	cdc.RegisterConcrete(Datagram{}, "cosmos-sdk/ibc/Datagram", nil)
	cdc.RegisterInterface((*Payload)(nil), nil)
}
