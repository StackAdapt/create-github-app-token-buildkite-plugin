# Create GitHub App Token

A Buildkite plugin to generate GitHub App installation tokens for use in your pipelines. The created token expires after 1 hour (https://docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app/generating-an-installation-access-token-for-a-github-app)

## Usage in Buildkite Pipeline

```yaml
steps:
  - label: 'Run with GitHub token'
    plugins:
      - ./.buildkite/create-github-app-token:
          app-id: ${GITHUB_APP_ID}
          private-key: ${GITHUB_PRIVATE_KEY}
          env-var-name: GITHUB_TOKEN
    command: |
      # Now you can use $GITHUB_TOKEN in your commands
      echo "Using GitHub token: $GITHUB_TOKEN"
```

## Plugin Parameters

- `app-id`: (Optional) The GitHub App ID (if not provided, will be retrieved from the environment variable)
- `private-key`: (Optional) The private key content for the GitHub App (if not provided, will be retrieved from the environment variable)
- `env-var-name`: (Required) Name of the environment variable to set (e.g., GITHUB_TOKEN)

## How It Works

1. The plugin authenticates as a GitHub App using the provided App ID and private key
2. It retrieves an installation token for the app (expires after 1 hour)
3. It sets the token as a Buildkite environment variable with the specified name
4. The token can then be used in subsequent steps in your pipeline:
   - As an environment variable in the current job

## Development

To modify the plugin:

1. Update the Go code in `main.go` and related files
2. Commit the changes
