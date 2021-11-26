# GitHub App Authentication for integration with GitHub

## Introduction

GitHub Apps are the officially recommended way to integrate with GitHub because of their support for granular permissions to access data. For more information see [About Apps](https://docs.github.com/en/developers/apps/getting-started-with-apps/about-apps)

The `github-app-auth` application is specifically designed to enable integration of third-party CI/CD systems with GitHub by generating a token that can be used to interact with the GitHub API available to GitHub Apps.
A list of endpoints available to GitHub Apps is documented [here](https://docs.github.com/en/rest/overview/endpoints-available-for-github-apps)

## Examples

### Retrieving a list of repositories with the GH CLI

The [GitHub CLI](https://cli.github.com/) allows for convenient access to GitHub from the command line.
We can retrieve a list of repositories the GitHub App has permission to access by invoking it with the `GITHUB_TOKEN` environment variable set to the installation token generated by `github-app-auth`.

```bash
GITHUB_TOKEN=$(github-app-auth <app-id> <private-key>) gh repo list
```

- `<app-id>` is the GitHub App ID
- `<private-key>` is the path to the GitHub App PEM encoded private key

### Uploading a SARIF file

The GitHub [documentation](https://docs.github.com/en/code-security/code-scanning/using-codeql-code-scanning-with-your-existing-ci-system/configuring-codeql-cli-in-your-ci-system#uploading-results-to-github) for using CodeQL in a CI system provides the following example for uploading results.

```bash
echo "$UPLOAD_TOKEN" | codeql github upload-results --repository=<repository-name> \
      --ref=<ref> --commit=<commit> --sarif=<file> \
      --github-auth-stdin
```

The `$UPLOAD_TOKEN` must be a token with the `security_events` scope as described in the CodeQL manual [here](https://codeql.github.com/docs/codeql-cli/manual/github-upload-results/).

With `github-app-auth` application that relies on a GitHub App to generate a token the example becomes.

```bash
github-app-auth <app-id> <private-key> | codeql github upload-results --repository=<repository-name> \
      --ref=<ref> --commit=<commit> --sarif=<file> \
      --github-auth-stdin
```

- `<app-id>` is the GitHub App ID
- `<private-key>` is the path to the GitHub App PEM encoded private key
