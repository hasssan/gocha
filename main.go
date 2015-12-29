// Package main handles the cli and configuration,
// delegates the actions to specialised packages.
package main

import (
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/jgautheron/gocha/bumper"
	"github.com/jgautheron/gocha/changelog"
	"github.com/jgautheron/gocha/config"
	"github.com/jgautheron/gocha/logger"
	"github.com/jgautheron/gocha/repository"
)

// IDEAS
// - one CHANGELOG per release
// - codenames only for major & minor
// - Makefile (make check, make test)

const (
	argLogLevel = "log-level"
	argRepoPath = "repo-path"

	// Git Signature
	argUserName  = "username"
	argUserEmail = "email"

	// Push settings
	argPushStrategy   = "push-strategy"
	argPushUsername   = "push-username"
	argPushPublicKey  = "push-public-key"
	argPushPrivateKey = "push-private-key"
	argPushPassphrase = "push-passphrase"

	// Changelog settings
	argAppName    = "app-name"
	argAppTag     = "tag"
	argOutputFile = "output"

	// Commands
	cmdBump              = "bump"
	cmdBumpMajor         = "major"
	cmdBumpMinor         = "minor"
	cmdBumpPatch         = "patch"
	cmdChangelog         = "changelog"
	cmdChangelogGenerate = "generate"
)

var (
	// Build vars
	// Do not set these manually! these variables
	// are meant to be set through ldflags
	buildTag  string
	buildDate string
)

func main() {
	app := cli.NewApp()
	app.Name = "gocha"
	app.Usage = "a tool to help you manage versions and changelogs"
	app.Author = "Jonathan Gautheron"
	app.Email = "jgautheron@neverblend.in"

	// These two variables are injected at build time with ldflags
	if buildTag != "" && buildDate != "" {
		app.Version = buildTag + " built on " + buildDate
	}

	// Global flags
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   argLogLevel,
			EnvVar: "LOG_LEVEL",
			Usage:  "log level: debug, info, warning|warn, error, fatal or panic",
		},
		cli.StringFlag{
			Name:   argRepoPath,
			Value:  "./",
			EnvVar: "REPO_PATH",
			Usage:  "path to the repository",
		},

		// Git Signature
		cli.StringFlag{
			Name:   argUserName,
			EnvVar: "USER_NAME",
			Usage:  "user name used for the git commands",
		},
		cli.StringFlag{
			Name:   argUserEmail,
			EnvVar: "USER_EMAIL",
			Usage:  "user email used for the git commands",
		},

		// Push settings
		cli.StringFlag{
			Name:   argPushStrategy,
			EnvVar: "PUSH_STRATEGY",
			Usage:  "push strategy: ssh-agent, ssh-key",
		},
		cli.StringFlag{
			Name:   argPushUsername,
			Value:  "git",
			EnvVar: "PUSH_USERNAME",
			Usage:  "push username, ex. [git]@mydomain.com...",
		},
		cli.StringFlag{
			Name:   argPushPublicKey,
			EnvVar: "PUSH_PUBLIC_KEY",
			Usage:  "path to the public key",
		},
		cli.StringFlag{
			Name:   argPushPrivateKey,
			EnvVar: "PUSH_PRIVATE_KEY",
			Usage:  "path to the private key",
		},
		cli.StringFlag{
			Name:   argPushPassphrase,
			Value:  "",
			EnvVar: "PUSH_PASSPHRASE",
			Usage:  "passphrase for the private key",
		},
	}

	app.Commands = []cli.Command{{
		Name:  cmdBump,
		Usage: "bump the current version number, major, minor or patch",
		Subcommands: []cli.Command{
			{
				Name:  cmdBumpMajor,
				Usage: "major version bump",
				Action: func(c *cli.Context) {
					initBump(c, cmdBumpMajor)
				},
			},
			{
				Name:  cmdBumpMinor,
				Usage: "minor version bump",
				Action: func(c *cli.Context) {
					initBump(c, cmdBumpMinor)
				},
			},
			{
				Name:  cmdBumpPatch,
				Usage: "patch version bump",
				Action: func(c *cli.Context) {
					initBump(c, cmdBumpPatch)
				},
			},
		},
	}, {
		Name:  cmdChangelog,
		Usage: "manipulate the changelog",
		Subcommands: []cli.Command{
			{
				Name:   cmdChangelogGenerate,
				Usage:  "generate the changelog",
				Action: initGenerate,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:   argAppName,
						EnvVar: "APP_NAME",
						Usage:  "the application name",
					},
					cli.StringFlag{
						Name:   argAppTag,
						EnvVar: "APP_TAG",
						Usage:  "generate the changelog from the given tag",
					},
					cli.StringFlag{
						Name:   argOutputFile,
						EnvVar: "OUTPUT_FILE",
						Usage:  "output file path",
					},
				},
			},
		},
	},
	}

	app.Run(os.Args)
}

func initialize(c *cli.Context) *repository.Repository {
	// Configure logging
	loglvl := config.GetCliOrConfigString(argLogLevel, c.GlobalString(argLogLevel))
	if len(loglvl) == 0 {
		loglvl = log.InfoLevel.String()
	}
	logger.SetLogLevel(loglvl)

	rp, err := repository.New(c.GlobalString(argRepoPath))
	if err != nil {
		log.Fatal(err)
	}

	var user *repository.User
	var push *repository.Push

	// Get the user name & email for git signatures
	un := config.GetCliOrConfigString(argUserName, c.GlobalString(argUserName))
	ue := config.GetCliOrConfigString(argUserEmail, c.GlobalString(argUserEmail))
	if len(un) == 0 || len(ue) == 0 {
		user = &repository.User{
			Name:  un,
			Email: ue,
		}
	} else {
		// If the username and email were not defined,
		// try to get them from the local git config
		user = repository.GetUserFromGitConfig()
		if user == nil {
			log.Fatal("The username and email are not defined")
		}
	}

	// Get the push settings
	push = &repository.Push{
		Strategy:   config.GetCliOrConfigString("push/strategy", c.GlobalString(argPushStrategy)),
		Username:   config.GetCliOrConfigString("push/username", c.GlobalString(argPushUsername)),
		PublicKey:  config.GetCliOrConfigString("push/public-key", c.GlobalString(argPushPublicKey)),
		PrivateKey: config.GetCliOrConfigString("push/private-key", c.GlobalString(argPushPrivateKey)),
		Passphrase: config.GetCliOrConfigString("push/passphrase", c.GlobalString(argPushPassphrase)),
	}

	creds := &repository.Credentials{
		User: user,
		Push: push,
	}
	rp.SetCredentials(creds)

	return rp
}

// initialize wraps the processor call and directly passes cli values.
func initBump(c *cli.Context, bmp string) {
	rp := initialize(c)
	bumper.Up(rp, bmp)
}

func getAppName(c *cli.Context) string {
	if len(c.String(argAppName)) != 0 {
		return c.String(argAppName)
	}

	path := c.GlobalString(argRepoPath)

	// Use the folder name as the application name
	idx := strings.LastIndex(path, "/")
	if idx != -1 {
		return path[idx+1:]
	}

	return path
}

func initGenerate(c *cli.Context) {
	outputFile := c.String(argOutputFile)
	if len(outputFile) == 0 {
		outputFile = c.GlobalString(argRepoPath)
	}

	rp := initialize(c)
	changelog.Generate(rp, c.String(argAppTag), getAppName(c), outputFile)
}
