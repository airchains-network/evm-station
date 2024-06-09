package junction

import (
	"context"
	"cosmossdk.io/log"
	"github.com/airchains-network/junction/x/junction/types"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"
	"os"
	"time"
)

func VerifyPod(podNumber uint64, ctx context.Context, jClient cosmosclient.Client, account cosmosaccount.Account, addr, stationId, previousMerkleRootHash, MerkleRootHash string, Proof []byte) (success bool) {

	podDetails := QueryPod(jClient, ctx, stationId, podNumber)
	if podDetails == nil {
		log.NewLogger(os.Stderr).Debug("Pod not submitted, can not verify")
		return false
	} else if podDetails.IsVerified == true {
		log.NewLogger(os.Stderr).Debug("Pod already verified")
		return true
	}
	verifyPodStruct := types.MsgVerifyPod{
		Creator:                addr,
		StationId:              stationId,
		PodNumber:              podNumber,
		MerkleRootHash:         MerkleRootHash,
		PreviousMerkleRootHash: previousMerkleRootHash,
		ZkProof:                Proof,
	}

	for {
		txRes, errTxRes := jClient.BroadcastTx(ctx, account, &verifyPodStruct)
		if errTxRes != nil {
			errTxResStr := errTxRes.Error()
			log.NewLogger(os.Stderr).Debug("Error in VerifyPod transaction", "error", errTxResStr)
			log.NewLogger(os.Stderr).Debug("Retrying VerifyPod transaction after 10 seconds..")
			time.Sleep(10 * time.Second)
		} else {
			log.NewLogger(os.Stderr).Info("Pod Verification Tx Success", "txHash", txRes.TxHash)
			return true
		}
	}
}
