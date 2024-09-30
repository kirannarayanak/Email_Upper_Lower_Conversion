package main

import (
	"fmt"
	"hasura-email-pipeline/internal/hasura"
	"log"
	"os"
	"strings"
)

func main() {
	// Hasura GraphQL endpoint and admin secret
	endpoint := "https://hasura.app.ayshei.com/v1/graphql"
	secret := "!qJm8YmN2Cu@Dc_uBJU6h2CXoCE_QjLBs4UwME3cN-"

	// Initialize the Hasura client
	client := hasura.NewHasuraClient(endpoint)

	// Fetch all emails
	emails, err := client.FetchEmails(secret)
	if err != nil {
		log.Fatalf("Failed to fetch emails: %v", err)
	}

	// Process emails
	for _, email := range emails {
		// Check if email is in uppercase
		if email == strings.ToUpper(email) {
			// Convert email to lowercase
			newEmail := strings.ToLower(email)

			// Update the email in Hasura
			err := client.UpdateEmail(secret, email, newEmail)
			if err != nil {
				log.Printf("Failed to update email %s: %v", email, err)
			} else {
				fmt.Printf("Successfully updated %s to %s\n", email, newEmail)
			}
		}
	}

	// Optional: Write updated emails to a file for verification
	err = os.WriteFile("updated_emails.txt", []byte(fmt.Sprintln(emails)), 0644)
	if err != nil {
		log.Fatalf("Failed to write emails to file: %v", err)
	}
}
