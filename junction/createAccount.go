package junction

import (
	"cosmossdk.io/log"
	"fmt"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosaccount"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

func CreateAccount(accountName, accountPath, addressPrefix string) (err error) {

	registry, err := cosmosaccount.New(cosmosaccount.WithHome(accountPath))
	if err != nil {
		return err
	}

	// create wallet

	account, mnemonic, err := registry.Create(accountName)
	if err != nil {
		return err
	}

	addr, err := account.Address(addressPrefix)
	if err != nil {
		return err
	}

	// save account address, mnemonic, creation time, address prefix in the file <accountPath>/<accountName>.yaml
	details := &AccountDetails{
		Address:       addr,
		Mnemonic:      mnemonic,
		CreationTime:  time.Now(),
		AddressPrefix: addressPrefix,
	}
	data, err := yaml.Marshal(details)
	if err != nil {
		return err
	}
	filename := fmt.Sprintf("%s/%s.yaml", accountPath, accountName)
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}
	log.NewLogger(os.Stderr).Info("Account details saved in file.", "filename", filename)
	log.NewLogger(os.Stderr).Info("Account created", "address", addr, "mnemonic", mnemonic)
	log.NewLogger(os.Stderr).Warn("Put balance in Account before starting tracks", "address", addr)
	return nil
}
