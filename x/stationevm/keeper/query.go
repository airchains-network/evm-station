package keeper

import (
	"station-evm/x/stationevm/types"
)

var _ types.QueryServer = Keeper{}
