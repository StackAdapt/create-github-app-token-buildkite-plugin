package secret

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"create-github-app-token/internal/common"
)

var ExecCommand = exec.Command

// GetEnvVar retrieves an environment variable with the following precedence:
// 1. Plugin prefix (BUILDKITE_PLUGIN_CREATE_GITHUB_APP_TOKEN_)
// 2. Direct environment variable
// 3. Buildkite agent secret (if secretFallback is true)
// 4. Default value (if provided)
// Returns the value and a boolean indicating if the value was found
func GetEnvVar(name string, defaultValue string, secretFallback bool) (string, bool) {
	// Try with plugin prefix first
	value, found := os.LookupEnv(common.PluginPrefix + name)
	if found {
		return value, true
	}

	// Try direct environment variable
	value, found = os.LookupEnv(name)
	if found {
		return value, true
	}

	// Try as a Buildkite secret if secretFallback is enabled
	if secretFallback {
		secretValue, err := getSecret(name)
		if err == nil && secretValue != "" {
			return secretValue, true
		}
	}

	// Return default value if provided
	if defaultValue != "" {
		return defaultValue, true
	}

	return "", false
}

// getSecret retrieves a secret using the buildkite-agent secret command
func getSecret(name string) (string, error) {
	cmd := ExecCommand("buildkite-agent", "secret", "get", name)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve secret: %v", err)
	}
	return strings.TrimSpace(string(output)), nil
}
