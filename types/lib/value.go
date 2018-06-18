package lib

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Value struct {
	store sdk.KVStore
	cdc   *codec.Codec
	key   []byte
}

func NewValue(store sdk.KVStore, cdc *codec.Codec, key []byte) Value {
	return Value{
		store: store,
		cdc:   cdc,
		key:   key,
	}
}

func (v Value) MustGet(ptr interface{}) {
	bz := v.store.Get(v.key)
	v.cdc.MustUnmarshalBinary(bz, ptr)
}

func (v Value) Get(ptr interface{}) bool {
	bz := v.store.Get(v.key)
	if bz == nil {
		return false
	}
	v.cdc.MustUnmarshalBinary(bz, ptr)
	return true
}

func (v Value) Has() bool {
	bz := v.store.Get(v.key)
	return bz != nil
}

func (v Value) Set(val interface{}) {
	v.store.Set(v.key, v.cdc.MustMarshalBinary(val))
}
