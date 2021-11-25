# GitHub App Authentication for CodeQL integration with GitHub

GitHub Apps are the officially recommended way to integrate with GitHub because of their support for granular permissions to access data. For more information see [About Apps](https://docs.github.com/en/developers/apps/getting-started-with-apps/about-apps)

The `github-app-auth` application is specifically designed to enable CodeQL integration with third-party CI/CD systems and generates a token that can be used to upload the results of a CodeQL analysis.

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
