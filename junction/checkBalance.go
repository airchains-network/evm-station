package junction

import (
	"context"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"
	"os"
)

func CheckBalance(client *cosmosclient.Client, accountAddress string) (bool, int64, error) {

	// check if client is connected or not

	ctx := context.Background()
	pageRequest := &query.PageRequest{} // Add this line to create a new PageRequest
	log.NewLogger(os.Stderr).Info("Checking balance", "accountAddress", accountAddress)

	balances, err := client.BankBalances(ctx, accountAddress, pageRequest) // Add pageRequest as the third argument
	if err != nil {
		log.NewLogger(os.Stderr).Error("Error querying bank balances:", "err", err)
		return false, 0, err
	}

	for _, balance := range balances {
		if balance.Denom == "amf" {
			return true, balance.Amount.Int64(), nil
		}
	}

	return false, 0, nil
}
