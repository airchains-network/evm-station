package tracks

//
//import (
//	"cosmossdk.io/log"
//	"github.com/airchains-network/evm-station/junction"
//	"github.com/airchains-network/evm-station/shared"
//	"github.com/airchains-network/evm-station/station"
//	"github.com/airchains-network/evm-station/tracks/tracksdb"
//	"github.com/airchains-network/evm-station/types"
//	"os"
//	"path/filepath"
//)
//
//func Start(configs *types.TracksConfigs, TracksDir string, homeDir string) {
//
//	// Initialise or Check database for Tracks
//	success := tracksdb.ConfigDb(TracksDir)
//	if !success {
//		log.NewLogger(os.Stderr).Error("Failed to initialize Tracks Database")
//		return
//	}
//
//	// get VRF private and public keys
//	VRFPrivateKey, err := junction.GetVRFPrivKey(homeDir)
//	if err != nil {
//		log.NewLogger(os.Stderr).Error("Failed to get VRF private key")
//		return
//	}
//	if VRFPrivateKey == "" {
//		log.NewLogger(os.Stderr).Error("VRF private key is empty")
//		return
//	}
//
//	VRFPublicKey, err := junction.GetVRFPubKey(homeDir)
//	if err != nil {
//		log.NewLogger(os.Stderr).Error("Failed to get VRF public key")
//		return
//	}
//	if VRFPublicKey == "" {
//		log.NewLogger(os.Stderr).Error("VRF public key is empty")
//		return
//	}
//
//	// Start remote signer (must start before node if running builtin).
//	log.NewLogger(os.Stderr).Info("Starting Tracks", "ChainId", configs.StationId, "JunctionRpc", configs.JunctionRpc, "JunctionKeyName", configs.JunctionKeyName)
//
//	TracksConfigs, err := getTracksData(cmd)
//	if err != nil {
//		log.NewLogger(os.Stderr).Error(err.Error())
//		return
//	}
//	accountName := TracksConfigs.JunctionKeyName
//	keyringDir := filepath.Join(homeDir, JunctionKeysFolder)
//	addr, err := junction.CheckIfAccountExists(accountName, keyringDir)
//	if err != nil {
//		log.NewLogger(os.Stderr).Error("Error in getting account address: " + err.Error())
//		return
//	}
//
//	jClient, jConnected := shared.GetJunctionClient()
//	if !jConnected {
//		log.NewLogger(os.Stderr).Error("Junction client not connected")
//		return
//	}
//
//	haveBalance, value, err := junction.CheckBalance(jClient, addr)
//	if err != nil {
//		log.NewLogger(os.Stderr).Error("Error in checking balance of"+addr, "Error", err.Error())
//		return
//	} else if !haveBalance {
//		log.NewLogger(os.Stderr).Warn("Not have balance in " + addr)
//		return
//	}
//
//	log.NewLogger(os.Stderr).Info("Junction Account Balance (in amf):", "account", addr, "balance", value)
//
//
//	for {
//		n := station.GetLatestPodNumber() // latest pod number"
//		log.NewLogger(os.Stderr).Info("latest pod", "pod_number", n)
//
//		junction.InitVRF()
//		//junction.ValidateVRF()
//
//		//junction.SubmitPod()
//		//junction.VerifyPod()
//
//		break
//	}
//}
