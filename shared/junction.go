package shared

import (
	"context"
	"cosmossdk.io/log"
	"fmt"
	"github.com/airchains-network/evm-station/junction"
	"github.com/airchains-network/evm-station/utils"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"
	"os"
)

var JunctionClient *cosmosclient.Client
var JunctionConnected bool

type JunctionConfigStruct struct {
	RPC           string
	AddressPrefix string
	HomeDir       string
}

func JunctionNewConfig(junctionRpc, keyHomePath string) *JunctionConfigStruct {
	return &JunctionConfigStruct{
		RPC:           junctionRpc,
		AddressPrefix: junction.JunctionAddressPrefix,
		HomeDir:       keyHomePath,
	}
}

func SetJunctionClient(j *JunctionConfigStruct) {

	gas := utils.GenerateRandomWithFavour(611, 1200, [2]int{612, 1000}, 0.7)

	gasFees := fmt.Sprintf("%damf", gas)

	options := []cosmosclient.Option{
		cosmosclient.WithNodeAddress(j.RPC),
		cosmosclient.WithAddressPrefix(j.AddressPrefix),
		cosmosclient.WithHome(j.HomeDir),
		cosmosclient.WithGas("auto"),
		cosmosclient.WithFees(gasFees),
	}

	ctx := context.Background()
	client, err := cosmosclient.New(ctx, options...)

	if err != nil {
		log.NewLogger(os.Stderr).Error(err.Error())
		JunctionConnected = false
	} else {
		JunctionClient = &client
		JunctionConnected = true
		log.NewLogger(os.Stderr).Info("Junction Client Connected Successfully")
	}
}

func GetJunctionClient() (*cosmosclient.Client, bool) {
	return JunctionClient, JunctionConnected
}
