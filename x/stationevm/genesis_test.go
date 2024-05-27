package stationevm_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "station-evm/testutil/keeper"
	"station-evm/testutil/nullify"
	"station-evm/x/stationevm"
	"station-evm/x/stationevm/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.StationevmKeeper(t)
	stationevm.InitGenesis(ctx, *k, genesisState)
	got := stationevm.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
