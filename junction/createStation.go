package junction

import (
	"context"
	"cosmossdk.io/log"
	"encoding/json"
	"github.com/airchains-network/evm-station/types"
	junctionTypes "github.com/airchains-network/junction/x/junction/types"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"
	"os"
)

func CreateStation(junctionClient cosmosclient.Client, keyringDir, accountName, addr string, extraArg junctionTypes.StationArg, stationId string, stationInfo types.StationInfo, verificationKey groth16.VerifyingKey, tracks []string) bool {

	registry, err := cosmosaccount.New(cosmosaccount.WithHome(keyringDir))
	if err != nil {
		log.NewLogger(os.Stderr).Error("error in getting registry" + err.Error())
		return false
	}
	account, err := registry.GetByName(accountName)
	if err != nil {
		log.NewLogger(os.Stderr).Error("error in getting account" + err.Error())
		return false
	}
	addr, err = CheckIfAccountExists(accountName, keyringDir)
	if err != nil {
		log.NewLogger(os.Stderr).Error("Error in getting account address: " + err.Error())
		return false
	}

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
		log.NewLogger(os.Stderr).Error("Error in broadcasting transaction", "Error", err.Error())
		return false
	}

	log.NewLogger(os.Stdout).Info("Transaction sent successfully", "txHash", txResp.TxHash, "stationId", stationId)

	// create VRF Keys
	vrfPrivateKey, vrfPublicKey := NewKeyPair()
	vrfPrivateKeyHex := vrfPrivateKey.String()
	vrfPublicKeyHex := vrfPublicKey.String()
	if vrfPrivateKeyHex != "" {
		SetVRFPrivKey(vrfPrivateKeyHex)
	} else {
		log.NewLogger(os.Stderr).Error("Error saving VRF private key")
		return false
	}
	if vrfPublicKeyHex != "" {
		SetVRFPubKey(vrfPublicKeyHex)
	} else {
		log.NewLogger(os.Stderr).Error("Error saving VRF public key")
		return false
	}
	log.NewLogger(os.Stderr).Info("Successfully Created VRF public and private Keys")
	return true
}
