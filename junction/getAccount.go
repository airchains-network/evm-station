package junction

import (
	"cosmossdk.io/log"
	"fmt"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosaccount"
	"gopkg.in/yaml.v2"
	"os"
)

func RevealAccount(accountPath, addressPrefix, accountName string) (err error) {

	registry, err := cosmosaccount.New(cosmosaccount.WithHome(accountPath))
	if err != nil {
		log.NewLogger(os.Stderr).Error("Folder aka Registery do not exits ", "err", err)
		return err
	}

	// get account by name
	account, err := registry.GetByName(accountName)
	if err != nil {
		log.NewLogger(os.Stderr).Error("Error getting account", "err", err)
		return err
	}

	// get address from account
	addr, err := account.Address(addressPrefix)
	if err != nil {
		log.NewLogger(os.Stderr).Error("Error getting account address via address prefix", "err", err)
		return err
	}

	// get public key from account
	pubKey, err := account.PubKey()
	if err != nil {
		return err
	}

	// get mnemonic from .yaml if exists
	// read <accountPath>/<accountName>.yaml file -> unmarshal it -> get mnemonic
	filename := fmt.Sprintf("%s/%s.yaml", accountPath, account.Name)
	data, err := os.ReadFile(filename)
	if err != nil {
		log.NewLogger(os.Stderr).Error("Error reading account details file", "err", err)
		return err
	}
	var details AccountDetails
	err = yaml.Unmarshal(data, &details)
	if err != nil {
		log.NewLogger(os.Stderr).Error("Error unmarshalling account details", "err", err)
		return err
	}

	log.NewLogger(os.Stderr).Info("Account Details for ", "accountName", accountName)
	fmt.Println("Account Name:", account.Name)
	fmt.Println("Account Address:", addr)
	fmt.Println("Account Public Key:", pubKey)
	fmt.Println("Account Mnemonic:", details.Mnemonic)

	return nil
}

func GetAccount(accountPath, addressPrefix, accountName string) (accountAddress string, err error) {
	registry, err := cosmosaccount.New(cosmosaccount.WithHome(accountPath))
	if err != nil {
		log.NewLogger(os.Stderr).Error("Folder aka Registery do not exits ", "err", err)
		return "", err
	}

	// get account by name
	account, err := registry.GetByName(accountName)
	if err != nil {
		log.NewLogger(os.Stderr).Error("Error getting account", "err", err)
		return "", err
	}

	// get address from account
	addr, err := account.Address(addressPrefix)
	if err != nil {
		log.NewLogger(os.Stderr).Error("Error getting account address via address prefix", "err", err)
		return "", err
	}

	return addr, nil
}
