package sequencer

import (
	"fmt"
	"github.com/airchains-network/evm-station/types"
)

func Start(configs *types.SequencerConfigs) {

	// Start remote signer (must start before node if running builtin).
	fmt.Println("ChainId: ", configs.StationId)
	fmt.Println("JunctionRpc: ", configs.JunctionRpc)
	fmt.Println("JunctionKeyName: ", configs.JunctionKeyName)

	// take data from 26657

}
