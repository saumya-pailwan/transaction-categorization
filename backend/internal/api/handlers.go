package api

import (
	"context"
	"encoding/json"
	"net/http"

	"autonomoustx/internal/categorization"
	"autonomoustx/internal/db"
	"autonomoustx/internal/plaid"
)

func CreateLinkTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Get user ID from auth context
	userID := "user_good"
	token, err := plaid.CreateLinkToken(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"link_token": token})
}

func ExchangePublicTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PublicToken string `json:"public_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	accessToken, err := plaid.ExchangePublicToken(req.PublicToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Trigger initial sync
	go func() {
		// Sync transactions
		transactions, _, _, err := plaid.SyncTransactions(accessToken, "")
		if err != nil {
			// Log error
			return
		}

		// Process and save transactions
		for _, t := range transactions {
			// 1. Categorize
			var merchantName string
			if val := t.MerchantName.Get(); val != nil {
				merchantName = *val
			}

			tx := categorization.Transaction{
				Description: t.Name,
				Amount:      float64(t.Amount),
				Merchant:    merchantName,
			}

			result, _ := categorization.Categorize(tx)

			// 2. Generate Embedding
			embedding, _ := categorization.GenerateEmbedding(t.Name)

			// 3. Save to DB
			_, err := db.Pool.Exec(context.Background(), `
				INSERT INTO transactions (plaid_id, amount, date, description, merchant_name, category, embedding, is_manual)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
				ON CONFLICT (plaid_id) DO NOTHING
			`, t.TransactionId, t.Amount, t.Date, t.Name, t.MerchantName, result.Category, embedding, false)

			if err != nil {
				// Log error
			}
		}
	}()

	json.NewEncoder(w).Encode(map[string]string{"access_token": accessToken})
}

func GetTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Pool.Query(r.Context(), "SELECT id, description, amount, category, merchant_name, date FROM transactions ORDER BY date DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var transactions []map[string]interface{}
	for rows.Next() {
		var id int
		var desc, cat, merch string
		var amount float64
		var date string
		if err := rows.Scan(&id, &desc, &amount, &cat, &merch, &date); err != nil {
			continue
		}
		transactions = append(transactions, map[string]interface{}{
			"id":            id,
			"description":   desc,
			"amount":        amount,
			"category":      cat,
			"merchant_name": merch,
			"date":          date,
		})
	}

	json.NewEncoder(w).Encode(transactions)
}
