package main

import (
	"context"
	"encoding/pem"
	"fmt"
	"log"
	"os/exec"
	"strconv"

	"github.com/golang-jwt/jwt/v5"

	"create-github-app-token/internal/github"
	"create-github-app-token/internal/secret"
)

func main() {
	// Get environment variables using the secret package
	appIDStr, found := secret.GetEnvVar("APP_ID", "", true)
	if !found {
		log.Fatalf("Required environment variable APP_ID not set")
	}

	privateKeyContent, found := secret.GetEnvVar("PRIVATE_KEY", "", true)
	if !found {
		log.Fatalf("Required environment variable PRIVATE_KEY not set")
	}

	// For ENV_VAR_NAME, provide a default value of "GITHUB_TOKEN"
	envVarName, _ := secret.GetEnvVar("ENV_VAR_NAME", "GITHUB_TOKEN", false)

	// Convert app ID to int64
	appIDInt, err := strconv.ParseInt(appIDStr, 10, 64)
	if err != nil {
		log.Fatalf("Invalid app ID: %v", err)
	}

	// Parse private key
	block, _ := pem.Decode([]byte(privateKeyContent))
	if block == nil {
		log.Fatalf("Failed to parse PEM block containing the private key")
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyContent))
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	// Create a context
	ctx := context.Background()

	// Create GitHub client with JWT authentication
	client, err := github.NewClientWithJWT(appIDInt, privateKey)
	if err != nil {
		log.Fatalf("Failed to create GitHub client: %v", err)
	}

	// Create an installation token
	installationToken, err := client.CreateInstallationToken(ctx)
	if err != nil {
		log.Fatalf("Failed to create installation token: %v", err)
	}

	// Set Buildkite environment variable
	fmt.Printf("Setting Buildkite environment variable %s\n", envVarName)
	err = setBuildkiteEnvVar(envVarName, installationToken)
	if err != nil {
		log.Fatalf("Failed to set Buildkite environment variable: %v", err)
	}

	fmt.Printf("Successfully set %s in Buildkite environment variable\n", envVarName)
}

// setBuildkiteEnvVar sets a Buildkite metadata value and also sets it as an environment variable
func setBuildkiteEnvVar(name, value string) error {
	// Set as environment variable for the current job using buildkite-agent env set
	// This ensures the value is available as an environment variable in the current job
	fmt.Printf("Setting as environment variable for current job\n")
	envCmd := exec.Command("buildkite-agent", "env", "set", fmt.Sprintf("%s=%s", name, value))
	envOutput, err := envCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set environment variable: %v\nOutput: %s", err, string(envOutput))
	}
	fmt.Printf("Successfully set environment variable: %s\n", string(envOutput))
	
	return nil
}
