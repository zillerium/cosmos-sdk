package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Sum total of all staking tokens
func (k Keeper) TotalSupply(ctx sdk.Context) sdk.Dec {
	totalSupply, err := k.bankKeeper.GetDenomSupply(ctx, k.GetParams(ctx).BondDenom)
	if err != nil {
		panic("staking token doesn't exist in bank keeper")
	}

	return sdk.NewDecFromInt(totalSupply.Amount)
}

// LooseTokens - tokens which are not bonded in a validator
func (k Keeper) LooseTokens(ctx sdk.Context) sdk.Dec {
	return k.TotalSupply(ctx).Sub(k.GetBondedTokens(ctx))
}

// get the bond ratio of the global state
func (k Keeper) BondedRatio(ctx sdk.Context) sdk.Dec {
	supply := k.TotalSupply(ctx)
	if supply.GT(sdk.ZeroDec()) {
		return k.GetBondedTokens(ctx).Quo(supply)
	}
	return sdk.ZeroDec()
}

func (k Keeper) increaseBondedTokens(ctx sdk.Context, amt sdk.Dec) sdk.Dec {
	newAmt := k.GetBondedTokens(ctx).Add(amt)
	k.SetBondedTokens(ctx, newAmt)
	return newAmt
}

func (k Keeper) decreaseBondedTokens(ctx sdk.Context, amt sdk.Dec) sdk.Dec {
	newAmt := k.GetBondedTokens(ctx).Sub(amt)
	if newAmt.LT(sdk.ZeroDec()) {
		panic(fmt.Sprintf("sanity check: bonded tokens negative: %v", newAmt))
	}
	k.SetBondedTokens(ctx, newAmt)
	return newAmt
}

func (k Keeper) mintStakingTokens(ctx sdk.Context, amt sdk.Int) {
	currSupply, err := k.bankKeeper.GetDenomSupply(ctx, k.GetParams(ctx).BondDenom)
	if err != nil {
		panic("staking token doesn't exist in bank keeper")
	}
	addAmountCoin := sdk.NewCoin(k.GetParams(ctx).BondDenom, amt)
	newTotalSupply := currSupply.Plus(addAmountCoin)
	k.bankKeeper.SetDenomSupply(ctx, newTotalSupply)
}

func (k Keeper) burnStakingTokens(ctx sdk.Context, amt sdk.Int) {
	currSupply, err := k.bankKeeper.GetDenomSupply(ctx, k.GetParams(ctx).BondDenom)
	if err != nil {
		panic("staking token doesn't exist in bank keeper")
	}
	subAmountCoin := sdk.NewCoin(k.GetParams(ctx).BondDenom, amt)
	newTotalSupply := currSupply.Minus(subAmountCoin)
	k.bankKeeper.SetDenomSupply(ctx, newTotalSupply)
}

// HumanReadableString returns a human readable string representation of a
// pool.
func (k Keeper) PoolHumanReadableString(ctx sdk.Context) string {
	resp := "Pool \n"
	resp += fmt.Sprintf("Loose Tokens: %s\n", k.LooseTokens(ctx))
	resp += fmt.Sprintf("Bonded Tokens: %s\n", k.GetBondedTokens(ctx))
	resp += fmt.Sprintf("Token Supply: %s\n", k.TotalSupply(ctx))
	resp += fmt.Sprintf("Bonded Ratio: %v\n", k.BondedRatio(ctx))
	return resp
}
