package main

import (
	"fmt"
	"log"
	"os"

	vault "github.com/hashicorp/vault/api"
)

// initVaultClient initializes a Vault client
func initVaultClient() (*vault.Client, error) {
	config := vault.DefaultConfig()
	config.Address = os.Getenv("VAULT_ADDR")

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Vault client: %v", err)
	}

	token := os.Getenv("VAULT_TOKEN")
	client.SetToken(token)

	return client, nil
}

// getSecret retrieves a secret from Vault
func getSecret(client *vault.Client, path string) (map[string]interface{}, error) {
	secret, err := client.Logical().Read(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret from Vault: %v", err)
	}

	if secret == nil {
		return nil, fmt.Errorf("no secret found at path: %s", path)
	}

	return secret.Data, nil
}

// secureStorageExample demonstrates secure storage of sensitive data using Vault
func secureStorageExample() {
	client, err := initVaultClient()
	if err != nil {
		log.Fatalf("Error initializing Vault client: %v", err)
	}

	secretPath := "secret/data/api_key"
	secretData, err := getSecret(client, secretPath)
	if err != nil {
		log.Fatalf("Error retrieving secret: %v", err)
	}

	apiKey, ok := secretData["api_key"].(string)
	if !ok {
		log.Fatalf("Error: API key not found in secret data")
	}

	fmt.Printf("Retrieved API key: %s\n", apiKey)
}

// accessControlExample demonstrates access controls to restrict access to sensitive data
func accessControlExample() {
	client, err := initVaultClient()
	if err != nil {
		log.Fatalf("Error initializing Vault client: %v", err)
	}

	secretPath := "secret/data/api_key"
	secretData, err := getSecret(client, secretPath)
	if err != nil {
		log.Fatalf("Error retrieving secret: %v", err)
	}

	apiKey, ok := secretData["api_key"].(string)
	if !ok {
		log.Fatalf("Error: API key not found in secret data")
	}

	// Implement access control logic here
	// For example, restrict access based on user roles or permissions

	fmt.Printf("Access granted to API key: %s\n", apiKey)
}
