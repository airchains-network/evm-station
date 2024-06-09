package junction

import (
	"context"
	"github.com/airchains-network/junction/x/junction/types"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"
)

func QueryLatestVerifiedBatch(client cosmosclient.Client, ctx context.Context, stationId string) uint64 {
	queryClient := types.NewQueryClient(client.Context())
	queryResp, err := queryClient.GetLatestVerifiedPodNumber(ctx, &types.QueryGetLatestVerifiedPodNumberRequest{StationId: stationId})
	if err != nil {
		return 0
	}
	return queryResp.PodNumber
}
func QueryPod(client cosmosclient.Client, ctx context.Context, stationId string, podNumber uint64) (pod *types.Pods) {
	queryClient := types.NewQueryClient(client.Context())
	queryResp, err := queryClient.GetPod(ctx, &types.QueryGetPodRequest{StationId: stationId, PodNumber: podNumber})
	if err != nil {
		//logs.Log.Error("Error fetching VRF: " + err.Error())
		return nil
	}

	return queryResp.Pod
}
