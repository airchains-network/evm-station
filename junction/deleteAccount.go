package junction

import (
	"github.com/ignite/cli/v28/ignite/pkg/cosmosaccount"
	"log"
	"os"
)

func DeleteAccount(keyringDir, addressPrefix, accountName string) (err error) {

	// check if wallet exists
	registry, err := cosmosaccount.New(cosmosaccount.WithHome(keyringDir))
	if err != nil {
		return err
	}

	// get account
	_, err = registry.GetByName(accountName)
	if err != nil {
		log.New(os.Stderr, "", 0).Println("Account does not exist or already deleted")
		return err
	}

	// create wallet
	err = registry.DeleteByName(accountName)
	if err != nil {
		log.New(os.Stderr, "", 0).Println("Error deleting account")
		return err
	}

	return nil
}
