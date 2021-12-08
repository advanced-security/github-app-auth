package main

import (
	"context"
	"io/ioutil"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/go-github/v40/github"
	"golang.org/x/oauth2"
)

func createClaims(appId string) *jwt.StandardClaims {
	return &jwt.StandardClaims{
		// Issued at time, 60 seconds in the past to allow for clock drift
		IssuedAt: time.Now().Unix() - 60,
		// Expiration time, 10 minute maximum
		ExpiresAt: time.Now().Unix() + (10 * 60),
		// GitHub App's identifier
		Issuer: appId,
	}
}

type GitHubApp struct {
	appId      string
	privateKey []byte
	client     *github.Client
	ctx        context.Context
}

func NewGitHubApp(appId string, privateKeyFilePath string) (*GitHubApp, error) {
	app := &GitHubApp{appId: appId}
	key, err := ioutil.ReadFile(privateKeyFilePath)
	if err != nil {
		return nil, err
	}
	app.privateKey = key

	claims := createClaims(app.appId)

	signingKey, err := jwt.ParseRSAPrivateKeyFromPEM(app.privateKey)
	if err != nil {
		return nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, *claims)
	signedToken, err := token.SignedString(signingKey)
	if err != nil {
		return nil, err
	}

	app.ctx = context.Background()
	staticTokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: signedToken})
	tokenClient := oauth2.NewClient(app.ctx, staticTokenSource)
	app.client = github.NewClient(tokenClient)

	return app, nil
}

func (app *GitHubApp) ListInstallations() ([]*github.Installation, error) {
	installations, _, err := app.client.Apps.ListInstallations(app.ctx, nil)

	if err != nil {
		return nil, err
	}

	return installations, nil
}

func (app *GitHubApp) CreateInstallationToken(installation *github.Installation) (*github.InstallationToken, error) {

	installationToken, _, err := app.client.Apps.CreateInstallationToken(app.ctx, *installation.ID, &github.InstallationTokenOptions{Permissions: installation.Permissions})
	if err != nil {
		return nil, err
	}

	return installationToken, nil
}

func (app *GitHubApp) CreateGitHubClient(ctx context.Context, installation *github.Installation) (*github.Client, error) {

	installationToken, err := app.CreateInstallationToken(installation)
	if err != nil {
		return nil, err
	}

	if ctx == nil {
		ctx = context.Background()
	}
	staticTokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: *installationToken.Token})
	tokenClient := oauth2.NewClient(ctx, staticTokenSource)
	client := github.NewClient(tokenClient)

	return client, nil
}
