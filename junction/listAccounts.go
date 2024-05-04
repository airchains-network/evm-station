package junction

import (
	"cosmossdk.io/log"
	"fmt"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosaccount"
	"gopkg.in/yaml.v2"
	"os"
)

func ListAccounts(accountPath, addressPrefix string) (err error) {

	registry, err := cosmosaccount.New(cosmosaccount.WithHome(accountPath))
	if err != nil {
		return err
	}

	// List all accounts
	accounts, err := registry.List()
	if err != nil {
		log.NewLogger(os.Stderr).Error("Error listing accounts", "err", err)
		return err
	}

	if len(accounts) == 0 {
		log.NewLogger(os.Stderr).Info("No accounts found in this path:", accountPath)
		return nil
	}

	// Print all accounts
	log.NewLogger(os.Stderr).Info("All Accounts list. Count: " + fmt.Sprintf("%v", len(accounts)))
	for _, account := range accounts {
		// get address from account
		addr, err := account.Address(addressPrefix)
		if err != nil {
			return err
		}

		// get public key from account
		pubKey, err := account.PubKey()
		if err != nil {
			return err
		}
		// convert  PubKeySecp256k1{02BE5CAFD9B3FD2CFB2C16DE38CA16C237A9C67763972F30DBB92CB60620C50A8E} to 02BE5CAFD9B3FD2CFB2C16DE38CA16C237A9C67763972F30DBB92CB60620C50A8E
		pubKey = pubKey[16 : len(pubKey)-1]

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

		fmt.Println("------------------------------------------------")
		fmt.Println("Account Name:", account.Name)
		fmt.Println("Account Address:", addr)
		fmt.Println("Account Public Key:", pubKey)
		fmt.Println("Account Mnemonic:", details.Mnemonic)
	}

	return nil
}
