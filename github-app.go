package main

import (
	"context"
	"errors"
	"io/ioutil"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/go-github/v40/github"
	"golang.org/x/oauth2"
)

const (
	ReadPermission  = "read"
	WritePermission = "write"
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
}

func NewGitHubApp(appId string, privateKeyFilePath string) (*GitHubApp, error) {
	app := &GitHubApp{appId: appId}
	key, err := ioutil.ReadFile(privateKeyFilePath)
	if err != nil {
		return nil, err
	}
	app.privateKey = key

	return app, nil
}

func (app *GitHubApp) CreateInstallationToken() (*github.InstallationToken, error) {
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

	ctx := context.Background()
	staticTokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: signedToken})
	tokenClient := oauth2.NewClient(ctx, staticTokenSource)
	client := github.NewClient(tokenClient)

	installations, _, err := client.Apps.ListInstallations(ctx, nil)
	if err != nil {
		return nil, err
	}

	if len(installations) != 1 {
		return nil, errors.New("The GitHub App has multiple installations when only one is expected!")
	}

	installation := installations[0]
	permission := WritePermission
	installationToken, _, err := client.Apps.CreateInstallationToken(ctx, *installation.ID, &github.InstallationTokenOptions{Permissions: &github.InstallationPermissions{SecurityEvents: &permission}})
	if err != nil {
		return nil, err
	}

	return installationToken, nil
}

func (app *GitHubApp) CreateGitHubClient(ctx context.Context) (*github.Client, error) {

	installationToken, err := app.CreateInstallationToken()
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
