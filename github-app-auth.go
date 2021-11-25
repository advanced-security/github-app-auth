package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
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

func main() {
	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) != 2 {
		fmt.Fprintf(os.Stderr, "Usage %s appId privateKey", os.Args[0])
		os.Exit(1)
	}

	appId := argsWithoutProg[0]
	keyFilePath := argsWithoutProg[1]

	claims := createClaims(appId)
	keyBytes, err := ioutil.ReadFile(keyFilePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	signingKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyBytes)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, *claims)
	signedToken, err := token.SignedString(signingKey)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ctx := context.Background()
	staticTokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: signedToken})
	tokenClient := oauth2.NewClient(ctx, staticTokenSource)
	client := github.NewClient(tokenClient)

	installations, _, err := client.Apps.ListInstallations(ctx, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if len(installations) != 1 {
		fmt.Fprintln(os.Stderr, "The GitHub App has multiple installations when only one is expected!")
		os.Exit(1)
	}

	installation := installations[0]
	permission := WritePermission
	installationToken, _, err := client.Apps.CreateInstallationToken(ctx, *installation.ID, &github.InstallationTokenOptions{Permissions: &github.InstallationPermissions{SecurityEvents: &permission}})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("%s", *installationToken.Token)
}
