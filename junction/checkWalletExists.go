package junction

import (
	"github.com/ignite/cli/v28/ignite/pkg/cosmosaccount"
)

func CheckIfAccountExists(accountName, accountPath string) (addr string, err error) {
	registry, err := cosmosaccount.New(cosmosaccount.WithHome(accountPath))
	if err != nil {
		return "", err
	}

	account, err := registry.GetByName(accountName)
	if err != nil {
		return "", err
	}

	addr, err = account.Address(addressPrefix)
	if err != nil {
		return "", err
	}

	return addr, nil
}

func GetCosmosAccount(accountName, accountPath string) (account cosmosaccount.Account, err error) {
	registry, err := cosmosaccount.New(cosmosaccount.WithHome(accountPath))
	if err != nil {
		return account, err
	}

	account, err = registry.GetByName(accountName)
	if err != nil {
		return account, err
	}
	return account, nil
}
