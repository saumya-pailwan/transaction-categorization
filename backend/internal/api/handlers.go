package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

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

	// Trigger initial sync (Synchronous for first load to avoid race conditions)
	log.Printf("Starting initial sync for token: %s...", accessToken[:10])

	// Sync transactions
	transactions, _, _, err := plaid.SyncTransactions(accessToken, "")
	if err != nil {
		log.Printf("Error syncing transactions: %v", err)
		// We still return success for the token exchange, but maybe with a warning?
		// For now, let's just log it.
	} else {
		log.Printf("Synced %d transactions from Plaid", len(transactions))

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

			result, err := categorization.Categorize(tx)
			if err != nil {
				log.Printf("Categorization error for %s: %v", t.Name, err)
			}

			// 2. Generate Embedding
			var embeddingParam interface{} = nil
			vec, err := categorization.GenerateEmbedding(t.Name)
			if err != nil {
				log.Printf("Embedding error for %s: %v", t.Name, err)
			} else {
				embeddingParam = vec
			}

			// 3. Save to DB
			// Use extraction merchantName, NOT t.MerchantName
			_, err = db.Pool.Exec(context.Background(), `
				INSERT INTO transactions (plaid_id, amount, date, description, merchant_name, category, embedding, is_manual)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
				ON CONFLICT (plaid_id) DO NOTHING
			`, t.TransactionId, t.Amount, t.Date, t.Name, merchantName, result.Category, embeddingParam, false)

			if err != nil {
				log.Printf("DB Insert error: %v", err)
			}
		}
	}

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
		var desc string
		var cat, merch *string // Handle NULLs
		var amount float64
		var date time.Time // Handle DATE type

		if err := rows.Scan(&id, &desc, &amount, &cat, &merch, &date); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		var catStr, merchStr string
		if cat != nil {
			catStr = *cat
		}
		if merch != nil {
			merchStr = *merch
		}

		transactions = append(transactions, map[string]interface{}{
			"id":            id,
			"description":   desc,
			"amount":        amount,
			"category":      catStr,
			"merchant_name": merchStr,
			"date":          date.Format("2006-01-02"),
		})
	}

	json.NewEncoder(w).Encode(transactions)
}
