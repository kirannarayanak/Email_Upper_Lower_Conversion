package hasura

import (
	"context"
	"log"

	"github.com/machinebox/graphql"
)

// HasuraClient represents the client used to interact with Hasura.
type HasuraClient struct {
	client *graphql.Client
}

// NewHasuraClient creates a new GraphQL client for Hasura.
func NewHasuraClient(endpoint string) *HasuraClient {
	return &HasuraClient{
		client: graphql.NewClient(endpoint),
	}
}

// FetchEmails fetches all emails from the customer table
func (hc *HasuraClient) FetchEmails(secret string) ([]string, error) {
	req := graphql.NewRequest(`
        query {
            customer {
                email
            }
        }
    `)

	// Add Hasura admin secret
	req.Header.Set("X-Hasura-Admin-Secret", secret)

	// Response structure
	var respData struct {
		Customer []struct {
			Email string `json:"email"`
		} `json:"customer"`
	}

	// Execute the request
	if err := hc.client.Run(context.Background(), req, &respData); err != nil {
		return nil, err
	}

	// Collect and return emails
	emails := []string{}
	for _, customer := range respData.Customer {
		emails = append(emails, customer.Email)
	}

	return emails, nil
}

// UpdateEmail updates a customer's email in Hasura.
func (hc *HasuraClient) UpdateEmail(secret, oldEmail, newEmail string) error {
	req := graphql.NewRequest(`
        mutation UpdateCustomerEmail($oldEmail: String!, $newEmail: String!) {
            update_customer(where: {email: {_eq: $oldEmail}}, _set: {email: $newEmail}) {
                affected_rows
            }
        }
    `)

	// Add variables to the request
	req.Var("oldEmail", oldEmail)
	req.Var("newEmail", newEmail)

	// Add Hasura admin secret
	req.Header.Set("X-Hasura-Admin-Secret", secret)

	// Execute the mutation
	var respData struct {
		UpdateCustomer struct {
			AffectedRows int `json:"affected_rows"`
		} `json:"update_customer"`
	}
	if err := hc.client.Run(context.Background(), req, &respData); err != nil {
		return err
	}

	log.Printf("Updated %d row(s)", respData.UpdateCustomer.AffectedRows)
	return nil
}
