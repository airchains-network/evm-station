package junction

import (
	"context"
	"cosmossdk.io/log"
	"github.com/airchains-network/evm-station/types"
	junctionTypes "github.com/airchains-network/junction/x/junction/types"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"
	"os"
	"strings"
	"time"
)

func ValidateVRF(podNumber uint64, ctx context.Context, jClient cosmosclient.Client, account cosmosaccount.Account, addr, stationId, privateKeyStr, publicKey string) (success bool) {
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
	msg := junctionTypes.MsgValidateVrf{
		Creator:      addr,
		StationId:    stationId,
		PodNumber:    podNumber,
		SerializedRc: serializedRC,
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
			if strings.Contains(errStr, VRFValidatedErrorContains) {
				log.NewLogger(os.Stderr).Debug("VRF already verified for this pod number")
				return true
			} else {
				log.NewLogger(os.Stderr).Error("Error in ValidateVRF transaction")
				log.NewLogger(os.Stderr).Debug("Retrying ValidateVRF transaction after 10 seconds..")
				time.Sleep(10 * time.Second)
			}
		} else {
			log.NewLogger(os.Stderr).Info("VRF Validated Successfully", "txHash", txRes.TxHash)
			return true
		}
	}
}
