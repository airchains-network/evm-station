package junction

import (
	"context"
	"cosmossdk.io/log"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/airchains-network/evm-station/types"
	junctionTypes "github.com/airchains-network/junction/x/junction/types"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/edwards25519"
	"os"
	"strings"
	"time"
)

func InitVRF(podNumber uint64, ctx context.Context, jClient cosmosclient.Client, account cosmosaccount.Account, addr, stationId, privateKeyStr, publicKey string) (success bool) {

	suite := edwards25519.NewBlakeSHA256Ed25519()
	privateKey, err := LoadHexPrivateKey(privateKeyStr)
	if err != nil {
		log.NewLogger(os.Stderr).Error("Error in loading private key: " + err.Error())
		return false
	}
	tracks := []string{addr}
	upperBond := uint64(len(tracks))
	rc := types.RequestCommitmentV2Plus{
		BlockNum:         1,
		StationId:        stationId,
		UpperBound:       upperBond,
		RequesterAddress: addr,
	}
	serializedRC, err := SerializeRequestCommitmentV2Plus(rc)
	if err != nil {
		log.NewLogger(os.Stderr).Error(err.Error())
		return false
	}
	proof, vrfOutput, err := GenerateVRFProof(suite, privateKey, serializedRC, int64(rc.BlockNum))
	if err != nil {
		log.NewLogger(os.Stderr).Error("Error generating unique proof: " + err.Error())
		return false
	}
	type ExtraArg struct {
		SerializedRc []byte `json:"serializedRc"`
		Proof        []byte `json:"proof"`
		VrfOutput    []byte `json:"vrfOutput"`
	}
	extraArg := ExtraArg{
		SerializedRc: serializedRC,
		Proof:        proof,
		VrfOutput:    vrfOutput,
	}
	extraArgsByte, err := json.Marshal(extraArg)
	if err != nil {
		log.NewLogger(os.Stderr).Error("Error marshaling extra arg: " + err.Error())
		return false
	}
	var defaultOccupancy uint64
	defaultOccupancy = 1
	msg := junctionTypes.MsgInitiateVrf{
		Creator:        addr,
		PodNumber:      podNumber,
		StationId:      stationId,
		Occupancy:      defaultOccupancy,
		CreatorsVrfKey: publicKey,
		ExtraArg:       extraArgsByte,
	}

	// check if this pod is behind the current pod at switchyard or not
	latestVerifiedBatch := QueryLatestVerifiedBatch(jClient, ctx, stationId)
	if latestVerifiedBatch+1 != podNumber {
		log.NewLogger(os.Stderr).Debug("Incorrect pod number")
		if latestVerifiedBatch+1 < podNumber {
			log.NewLogger(os.Stderr).Debug("Rollback required")
			return false
		} else if latestVerifiedBatch+1 > podNumber {
			log.NewLogger(os.Stderr).Debug("Pod number at Switchyard is ahead of the current pod number")
			return true
		}
	}

	for {
		txRes, errTxRes := jClient.BroadcastTx(ctx, account, &msg)
		if errTxRes != nil {
			errStr := errTxRes.Error()
			if strings.Contains(errStr, VRFInitiatedErrorContains) {
				log.NewLogger(os.Stderr).Debug("VRF already initiated for this pod number")
				return true
			} else {
				log.NewLogger(os.Stderr).Error("Error in InitVRF transaction" + errStr)
				log.NewLogger(os.Stderr).Debug("Retrying InitVRF transaction after 10 seconds..")
				time.Sleep(10 * time.Second)
			}
		} else {
			log.NewLogger(os.Stderr).Info("VRF Initiated Successfully", "txHash", txRes.TxHash)
			return true
		}
	}
}

func LoadHexPrivateKey(hexPrivateKey string) (privateKey kyber.Scalar, err error) {
	privateKeyBytes, err := hex.DecodeString(hexPrivateKey)
	if err != nil {
		fmt.Printf("Error decoding private key: %v\n", err)
		return nil, err
	}
	suite := edwards25519.NewBlakeSHA256Ed25519()
	privateKey = suite.Scalar().SetBytes(privateKeyBytes)
	return privateKey, nil
}
