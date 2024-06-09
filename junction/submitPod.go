package junction

import (
	"context"
	"cosmossdk.io/log"
	"fmt"
	"github.com/airchains-network/junction/x/junction/types"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"
	"os"
	"time"
)

func SubmitPod(podNumber uint64, ctx context.Context, jClient cosmosclient.Client, account cosmosaccount.Account, addr, stationId, previousMerkleRootHash, MerkleRootHash string, Witness []byte, Proof []byte) (success bool) {

	var LatestPodStatusHashStr string
	LatestPodStatusHashStr = MerkleRootHash

	// previous pod hash
	var PreviousPodStatusHashStr string
	PreviousPodStatusHashStr = previousMerkleRootHash

	// get witness
	witnessByte := Witness

	unixTime := time.Now().Unix()
	currentTime := fmt.Sprintf("%d", unixTime)

	msg := types.MsgSubmitPod{
		Creator:                addr,
		StationId:              stationId,
		PodNumber:              podNumber,
		MerkleRootHash:         LatestPodStatusHashStr,
		PreviousMerkleRootHash: PreviousPodStatusHashStr, // bytes.NewBuffer(pMrh).String(), // PreviousPodStatusHashStr,
		PublicWitness:          witnessByte,
		Timestamp:              currentTime,
	}
	// check if pod is already submitted
	podDetails := QueryPod(jClient, ctx, stationId, podNumber)
	if podDetails != nil {
		log.NewLogger(os.Stderr).Info("Pod already submitted")
		return true
	}

	for {
		txRes, errTxRes := jClient.BroadcastTx(ctx, account, &msg)
		if errTxRes != nil {
			log.NewLogger(os.Stderr).Debug("Error in SubmitPod Transaction", "error", errTxRes.Error())
			log.NewLogger(os.Stderr).Debug("Retrying SubmitPod transaction after 10 seconds..")
			time.Sleep(10 * time.Second)
		} else {
			log.NewLogger(os.Stderr).Info("Pod Submit Successfully", "txHash", txRes.TxHash)
			return true
		}
	}
}
