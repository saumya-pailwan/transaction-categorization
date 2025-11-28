package plaid

import (
	"context"
	"fmt"

	"github.com/plaid/plaid-go/plaid"
)

func CreateLinkToken(userID string) (string, error) {
	ctx := context.Background()
	
	user := plaid.LinkTokenCreateRequestUser{
		ClientUserId: userID,
	}

	request := plaid.NewLinkTokenCreateRequest(
		"Tx Agent",
		"en",
		[]plaid.CountryCode{plaid.COUNTRYCODE_US},
		user,
	)
	
	request.SetProducts([]plaid.Products{plaid.PRODUCTS_TRANSACTIONS})

	resp, _, err := Client.PlaidApi.LinkTokenCreate(ctx).LinkTokenCreateRequest(*request).Execute()
	if err != nil {
		return "", fmt.Errorf("failed to create link token: %w", err)
	}

	return resp.GetLinkToken(), nil
}

func ExchangePublicToken(publicToken string) (string, error) {
	ctx := context.Background()
	
	request := plaid.NewItemPublicTokenExchangeRequest(publicToken)
	
	resp, _, err := Client.PlaidApi.ItemPublicTokenExchange(ctx).ItemPublicTokenExchangeRequest(*request).Execute()
	if err != nil {
		return "", fmt.Errorf("failed to exchange public token: %w", err)
	}

	return resp.GetAccessToken(), nil
}

func SyncTransactions(accessToken string, cursor string) ([]plaid.Transaction, string, bool, error) {
	ctx := context.Background()
	
	request := plaid.NewTransactionsSyncRequest(accessToken)
	if cursor != "" {
		request.SetCursor(cursor)
	}
	
	// Fetch up to 100 transactions
	count := int32(100)
	request.SetCount(count)

	resp, _, err := Client.PlaidApi.TransactionsSync(ctx).TransactionsSyncRequest(*request).Execute()
	if err != nil {
		return nil, "", false, fmt.Errorf("failed to sync transactions: %w", err)
	}

	return resp.GetAdded(), resp.GetNextCursor(), resp.GetHasMore(), nil
}
