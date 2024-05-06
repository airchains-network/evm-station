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

//
//func (config *Config) Seal() *Config {
//	config.mtx.Lock()
//	defer config.mtx.Unlock()
//
//	if config.sealed {
//		return config
//	}
//
//	config.sealed = true
//	close(config.Sealedch)
//
//	return config
//}
//
//func (config *Config) GetJunctionClient() *JunctionClient {
//	config.mtx.RLock()
//	defer config.mtx.RUnlock()
//	return config.junctionClient
//}
//
//var once sync.Once
//var instance *Config
//
//// This function is responsible for creating the instance and initializing the
//// junctionClient, which it does only once in a lifetime.
//func InitConfigAndClient() {
//	once.Do(func() {
//		instance = &Config{
//			Sealedch: make(chan struct{}),
//		}
//		// Here you initialize junctionClient.
//		j := NewJunction()
//		client := GetClient(j)
//		//junctionCli := *JunctionClient
//		instance.SetJunctionClient(&client)
//		instance.Seal()
//	})
//}
//
//// set junction client
//func (config *Config) SetJunctionClient(junctionClient *JunctionClient) {
//	config.mtx.Lock()
//	defer config.mtx.Unlock()
//	config.junctionClient = junctionClient
//}
//
//// This function is responsible for retrieving the previously initialized instance.
//func GetConfigInstance() *Config {
//	if instance == nil {
//		panic("Config was not initialized. Call InitConfigAndClient() first.")
//	}
//	return instance
//}
