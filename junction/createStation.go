package junction

import (
	"context"
	"cosmossdk.io/log"
	"encoding/json"
	junctionTypes "github.com/airchains-network/evm-station/junction/types"
	"github.com/airchains-network/evm-station/types"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"
	"os"
)

func CreateStation(junctionClient cosmosclient.Client, account cosmosaccount.Account, addr string, extraArg junctionTypes.StationArg, stationId string, stationInfo types.StationInfo, verificationKey groth16.VerifyingKey, addressPrefix string, tracks []string) bool {

	stationJsonBytes, err := json.Marshal(stationInfo)
	if err != nil {
		log.NewLogger(os.Stderr).Error("Error marshaling to JSON: " + err.Error())
		return false
	}
	stationInfoStr := string(stationJsonBytes)

	verificationKeyByte, err := json.Marshal(verificationKey)
	if err != nil {
		log.NewLogger(os.Stderr).Error("Failed to unmarshal Verification key" + err.Error())
		return false
	}

	extraArgBytes, err := json.Marshal(extraArg)
	if err != nil {
		log.NewLogger(os.Stderr).Error("Error marshalling extra arg")
		return false
	}

	var tracksVotingPower []uint64
	totalPower := uint64(100)
	numTracks := len(tracks)
	equalShare := totalPower / uint64(numTracks)
	remainder := totalPower % uint64(numTracks)
	for i := 0; i < numTracks; i++ {
		if remainder > 0 {
			tracksVotingPower = append(tracksVotingPower, equalShare+1)
			remainder-- // Decrement the remainder until it's 0
		} else {
			tracksVotingPower = append(tracksVotingPower, equalShare)
		}
	}

	ctx := context.Background()

	stationData := junctionTypes.MsgInitStation{
		Creator:           addr,
		Tracks:            tracks,
		VerificationKey:   verificationKeyByte,
		StationId:         stationId,
		StationInfo:       stationInfoStr,
		TracksVotingPower: tracksVotingPower,
		ExtraArg:          extraArgBytes,
	}

	txResp, err := junctionClient.BroadcastTx(ctx, account, &stationData)
	if err != nil {
		log.NewLogger(os.Stderr).Error("Error in broadcasting transaction" + err.Error())
		return false
	}
	log.NewLogger(os.Stdout).Info("Transaction sent successfully", "txHash", txResp.TxHash)

	return true

	//
	//	return true
	//	//timestamp := time.Now().String()
	//	//successGenesis := CreateGenesisJson(stationInfo, verificationKey, stationId, tracks, tracksVotingPower, txResp.TxHash, timestamp, extraArg, newTempAddr)
	//	//if !successGenesis {
	//	//	return false
	//	//}
	//	//
	//	//// create VRF Keys
	//	//vrfPrivateKey, vrfPublicKey := NewKeyPair()
	//	//vrfPrivateKeyHex := vrfPrivateKey.String()
	//	//vrfPublicKeyHex := vrfPublicKey.String()
	//	//if vrfPrivateKeyHex != "" {
	//	//	SetVRFPrivKey(vrfPrivateKeyHex)
	//	//} else {
	//	//	logs.Log.Error("Error saving VRF private key")
	//	//	return false
	//	//}
	//	//if vrfPublicKeyHex != "" {
	//	//	SetVRFPubKey(vrfPublicKeyHex)
	//	//} else {
	//	//	logs.Log.Error("Error saving VRF public key")
	//	//	return false
	//	//}
	//	//logs.Log.Info("Successfully Created VRF public and private Keys")
	//	//
	//	//homeDir, err := os.UserHomeDir()
	//	//if err != nil {
	//	//	logs.Log.Error("Error in getting home dir path: " + err.Error())
	//	//	return false
	//	//}
	//	//
	//	//ConfigFilePath := filepath.Join(homeDir, config.DefaultTracksDir, config.DefaultConfigDir, config.DefaultConfigFileName)
	//	//bytes, err := os.ReadFile(ConfigFilePath)
	//	//if err != nil {
	//	//	logs.Log.Error("Error reading sequencer.toml")
	//	//	return false
	//	//
	//	//}
	//	//
	//	//var conf config.Config // JunctionConfig
	//	//err = toml.Unmarshal(bytes, &conf)
	//	//if err != nil {
	//	//	logs.Log.Error("error in unmarshling file")
	//	//	return false
	//	//}
	//	//
	//	//// Update the values
	//	//conf.Junction.JunctionRPC = jsonRPC
	//	//conf.Junction.JunctionAPI = ""
	//	//conf.Junction.StationId = stationId
	//	//conf.Junction.VRFPrivateKey = vrfPrivateKeyHex
	//	//conf.Junction.VRFPublicKey = vrfPublicKeyHex
	//	//conf.Junction.AddressPrefix = "air"
	//	//conf.Junction.AccountPath = accountPath
	//	//conf.Junction.AccountName = accountName
	//	//conf.Junction.Tracks = tracks
	//	//
	//	//// Marshal the struct to TOML
	//	//f, err := os.Create(ConfigFilePath)
	//	//if err != nil {
	//	//	logs.Log.Error("Error creating file")
	//	//	return false
	//	//}
	//	//defer f.Close()
	//	//newData := toml.NewEncoder(f)
	//	//if err := newData.Encode(conf); err != nil {
	//	//}
	//	//
	//	//return true
	//
	//}
}
