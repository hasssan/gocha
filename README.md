# gocha [![Circle CI](https://circleci.com/gh/jgautheron/gocha.svg?style=svg)](https://circleci.com/gh/jgautheron/gocha) [![GoDoc](https://godoc.org/github.com/jgautheron/gocha?status.png)](https://godoc.org/github.com/jgautheron/gocha)

Gocha is an opinionated tool written in Golang that aim to help teams deal with semver versioning and changelogs.
Changelogs are generated from commit messages between tags, it supposes that you are using the [AngularJS Git Commit Message Conventions](about-angularjs-...).

// Note about lightweight vs annotated tags
// gocha will create annotated tags
// badges: LICENSE, codeship, godoc, coveralls

## Getting started
You can directly use the binary (see [Build](build)) or download the project.

```go
go get -u github.com/jgautheron/gocha
```

### Configuration
In order to keep simple the usage of the tool in command line, you can add a configuration file named `.gocha.yaml` at the root of your home folder.

```yaml
# ~/.gocha.yaml
log-level: debug

# used for signing Git operations
username: jgautheron
email: foo@bar.com

# push details
push:
  strategy: ssh-key # ssh-key or ssh-agent
  username: git # in most cases it's "git", used for pushing git@domain.com...
  public-key: ~/.ssh/id_rsa.pub
  private-key: ~/.ssh/id_rsa
  passphrase: 123
```

## Commands

### Global

`gocha` aims at simplifying your development process, that's why it will also need `Git` settings and permissions to handle transparently the boring stuff for you.

```
NAME:
   gocha - a tool to help you manage versions and changelogs

USAGE:
   gocha [global options] command [command options] [arguments...]
   
VERSION:
   1.0.0
   
AUTHOR(S):
   Jonathan Gautheron <jgautheron@neverblend.in> 
   
COMMANDS:
   bump     bump the current version number, major, minor or patch
   changelog    manipulate the changelog
   help, h  Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --log-level      log level: debug, info, warning|warn, error, fatal or panic [$LOG_LEVEL]
   --repo-path "./" path to the repository [$REPO_PATH]
   --username       user name used for the git commands [$USER_NAME]
   --email      user email used for the git commands [$USER_EMAIL]
   --push-strategy  push strategy: ssh-agent, ssh-key [$PUSH_STRATEGY]
   --push-username  push username, ex. [git]@mydomain.com... [$PUSH_USERNAME]
   --push-public-key    path to the public key [$PUSH_PUBLIC_KEY]
   --push-private-key   path to the private key [$PUSH_PRIVATE_KEY]
   --push-passphrase    passphrase for the private key [$PUSH_PASSPHRASE]
   --help, -h       show help
   --version, -v    print the version
```

#### `--username` and `--email`
Both are required for signing tags and commits

#### `--push*`
These options are required for pushing changes

Tip: using the configuration file will help you keep your command lines readable and simplify the usage of this tool.

#### Push informations

There are currently two push strategies available.
1. Using the SSH agent, the simplest way and recommended for OSX if you are using your keychain for storing credentials.
2. Using a SSH key, then you will have to pass the public key, private key and passphrase.

### `bump`

Bumps the version number based on the latest tag, then automatically pushes it. A codename is automatically generated.

```
NAME:
   gocha bump - bump the current version number, major, minor or patch

USAGE:
   gocha bump command [command options] [arguments...]

COMMANDS:
   major    major version bump
   minor    minor version bump
   patch    patch version bump
   help, h  Shows a list of commands or help for one command
   
OPTIONS:
   --help, -h   show help
```

### `changelog`

## Build

The binaries are downloadable in the [Github releases page](https://github.com/jgautheron/gocha/releases).
To generate a new binary, simply launch `make` at the root of the project.

### System compatibility
OS               | Status
---------------- | ------
OSX x86_64       | Supported, tested
Linux x86        | Supported, tested
Linux x86_64     | Supported, tested
Linux ARMv5      | Supported, tested
Linux ARMv7      | Supported, tested
Windows x86      | Supported
Windows x86_64   | Supported

## About semver
Semver stands for Semantic Versioning using the MAJOR.MINOR.PATCH notation, for more info: http://semver.org/.  
`gocha` doesn't necessarily want to lock you up with this type of versioning, if you'd like to use another semantic, create an issue or contribute!

## About AngularJS Git Commit Message Conventions
The AngularJS conventions are simple yet advanced, the format is previsible and easy to parse. The `scope` fits for many languages, ex. in Golang that would be packages. [Check the specification.](https://docs.google.com/document/d/1QrDFcIiPjSLDn3EL15IJygNPiHORgU1_OOAqWjiDU5Y/edit)  
if you'd like to use another set of conventions, create an issue or contribute!

## Contributing
Contributions are encouraged. Instructions are documented in [CONTRIBUTING.md](https://github.com/jgautheron/gocha/blob/master/CONTRIBUTING.md).

## Alternatives
- https://github.com/rafinskipg/git-changelog

## Credits
- https://github.com/sindresorhus/semver-regex
- https://github.com/rafinskipg/git-changelog

## Author
Jonathan Gautheron - jgautheron [A-T] neverblend.in  
https://twitter.com/jgautheron

## License
GPLv2