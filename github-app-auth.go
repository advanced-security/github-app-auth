package main

import (
	"fmt"
	"os"
)

func main() {
	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) != 2 {
		fmt.Fprintf(os.Stderr, "Usage %s appId privateKey\n", os.Args[0])
		os.Exit(1)
	}

	appId := argsWithoutProg[0]
	keyFilePath := argsWithoutProg[1]

	app, err := NewGitHubApp(appId, keyFilePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	installationToken, err := app.CreateInstallationToken()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("%s", *installationToken.Token)
}
