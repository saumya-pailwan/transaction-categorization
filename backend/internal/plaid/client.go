package plaid

import (
	"os"

	"github.com/plaid/plaid-go/plaid"
)

var Client *plaid.APIClient

func Init() {
	clientID := os.Getenv("PLAID_CLIENT_ID")
	secret := os.Getenv("PLAID_SECRET")
	env := os.Getenv("PLAID_ENV")

	if env == "" {
		env = "sandbox"
	}

	configuration := plaid.NewConfiguration()
	configuration.AddDefaultHeader("PLAID-CLIENT-ID", clientID)
	configuration.AddDefaultHeader("PLAID-SECRET", secret)
	
	switch env {
	case "production":
		configuration.UseEnvironment(plaid.Production)
	case "development":
		configuration.UseEnvironment(plaid.Development)
	case "sandbox":
		configuration.UseEnvironment(plaid.Sandbox)
	default:
		configuration.UseEnvironment(plaid.Sandbox)
	}

	Client = plaid.NewAPIClient(configuration)
}
