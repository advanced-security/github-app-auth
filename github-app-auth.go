package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/google/go-github/v40/github"
)

func main() {

	installationId := flag.Int("install-id", -1, "Installation ID identifying the installation to use for authentication. Defaults to the last installation.")
	owner := flag.String("owner", "", "Organization or user with an installation to use for authentication.")

	flag.Parse()

	args := flag.Args()

	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage %s appId privateKey\n", os.Args[0])
		os.Exit(1)
	}

	appId := args[0]
	keyFilePath := args[1]

	if *installationId != -1 && *owner != "" {
		fmt.Fprintln(os.Stderr, "Both an installation id and owner is specified which is not supported!")
		os.Exit(1)
	}

	app, err := NewGitHubApp(appId, keyFilePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	installations, err := app.ListInstallations()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if len(installations) == 0 {
		fmt.Fprintln(os.Stderr, "The GitHub App has no installations! For more information on how to install a GitHub App see: https://docs.github.com/en/developers/apps/managing-github-apps/installing-github-apps")
		os.Exit(1)
	}

	var installation *github.Installation
	if *installationId == -1 && *owner == "" {
		installation = installations[len(installations)-1]
	} else {
		for _, e := range installations {
			if *e.ID == int64(*installationId) {
				installation = e
				break
			}

			if *e.Account.Login == *owner {
				installation = e
				break
			}
		}

		if installation == nil {
			var errorMsg string
			if *installationId != -1 {
				errorMsg = fmt.Sprintf("Unable to find a suitable installation with id %d", *installationId)
			}
			if *owner != "" {
				errorMsg = fmt.Sprintf("Unable to find a suitable installation for owner %s", *owner)
			}

			fmt.Println(errorMsg)
			os.Exit(1)
		}
	}

	installationToken, err := app.CreateInstallationToken(installation)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("%s", *installationToken.Token)
}
