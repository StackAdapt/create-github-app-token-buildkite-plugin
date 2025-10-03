package github

import (
	"context"
	"crypto/rsa"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-github/v60/github"
)

// Client represents a GitHub client with app authentication capabilities
type Client struct {
	client *github.Client
	appID  int64
}

// NewClientWithJWT creates a new GitHub client authenticated with a JWT token
func NewClientWithJWT(appID int64, privateKey *rsa.PrivateKey) (*Client, error) {
	// Generate JWT token
	token, err := generateJWT(appID, privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT: %v", err)
	}

	// Create GitHub client with JWT authentication
	client := github.NewClient(nil).WithAuthToken(token)

	return &Client{
		client: client,
		appID:  appID,
	}, nil
}

// generateJWT creates a JWT token for GitHub App authentication
func generateJWT(appID int64, privateKey *rsa.PrivateKey) (string, error) {
	// JWT expiration time: 10 minutes from now
	now := time.Now()
	// Set issuedAt to 60 seconds in the past to prevent clock drift issues
	issuedAt := now.Add(-60 * time.Second)
	expirationTime := now.Add(10 * time.Minute)

	// Create the JWT claims
	claims := jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(issuedAt),
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		Issuer:    strconv.FormatInt(appID, 10),
	}

	// Create the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	
	// Sign the token with the private key
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %v", err)
	}

	return signedToken, nil
}

// CreateInstallationToken creates an installation token for the GitHub App
func (c *Client) CreateInstallationToken(ctx context.Context) (string, error) {
	// Get the app's installations
	installations, _, err := c.client.Apps.ListInstallations(ctx, &github.ListOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to list installations: %v", err)
	}

	if len(installations) == 0 {
		return "", fmt.Errorf("no installations found for app ID %d", c.appID)
	}

	// Use the first installation (you might want to add logic to find the correct one)
	installationID := installations[0].GetID()

	// Create an installation token
	token, _, err := c.client.Apps.CreateInstallationToken(ctx, installationID, &github.InstallationTokenOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to create installation token: %v", err)
	}

	return token.GetToken(), nil
}
