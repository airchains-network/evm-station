package tracks

import (
	"cosmossdk.io/log"
	"fmt"
	"github.com/airchains-network/evm-station/junction"
	"github.com/airchains-network/evm-station/tracks/tracksdb"
	"github.com/airchains-network/evm-station/types"
	"os"
)

func Start(configs *types.TracksConfigs, TracksDir string, homeDir string) {

	// Initialise or Check database for Tracks
	success := tracksdb.ConfigDb(TracksDir)
	if !success {
		log.NewLogger(os.Stderr).Error("Failed to initialize Tracks Database")
		return
	}

	// get VRF private and public keys
	VRFPrivateKey, err := junction.GetVRFPrivKey(homeDir)
	if err != nil {
		log.NewLogger(os.Stderr).Error("Failed to get VRF private key")
		return
	}
	if VRFPrivateKey == "" {
		log.NewLogger(os.Stderr).Error("VRF private key is empty")
		return
	}

	VRFPublicKey, err := junction.GetVRFPubKey(homeDir)
	if err != nil {
		log.NewLogger(os.Stderr).Error("Failed to get VRF public key")
		return
	}
	if VRFPublicKey == "" {
		log.NewLogger(os.Stderr).Error("VRF public key is empty")
		return
	}

	// Start remote signer (must start before node if running builtin).
	fmt.Println("ChainId: ", configs.StationId)
	fmt.Println("JunctionRpc: ", configs.JunctionRpc)
	fmt.Println("JunctionKeyName: ", configs.JunctionKeyName)

	// take data from 26657
}
