package shared

import (
	"context"
	"cosmossdk.io/log"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"
	"os"
)

var JunctionClient *cosmosclient.Client
var JunctionConnected bool

type JunctionConfigStruct struct {
	RPC            string
	KeyringDir     string
	KeyringBackend string
}

func JunctionNewConfig(junctionRpc, keyringPath string) *JunctionConfigStruct {
	return &JunctionConfigStruct{
		RPC:            junctionRpc,
		KeyringDir:     keyringPath,
		KeyringBackend: "os",
	}
}

func (j JunctionConfigStruct) JunctionConfigModify(rpc, keyringDir, keyringBackend string) {
	if rpc != "" {
		j.RPC = rpc
	}

	if keyringDir != "" {
		j.KeyringDir = keyringDir
	}

	if keyringBackend != "" {
		j.KeyringBackend = keyringBackend
	}
}

func SetJunctionClient(j *JunctionConfigStruct) {
	options := []cosmosclient.Option{
		cosmosclient.WithNodeAddress(j.RPC),
		cosmosclient.WithKeyringDir(j.KeyringDir),
		cosmosclient.WithKeyringBackend(cosmosaccount.KeyringBackend(j.KeyringBackend)),
	}

	ctx := context.Background()
	client, err := cosmosclient.New(ctx, options...)
	if err != nil {
		//log.NewLogger(os.Stderr).Error("Junction Client Connection Error", "Error", err)
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
