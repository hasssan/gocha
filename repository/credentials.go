// Package repository wraps and simplifies the libgit2 bindings
// exposed in the git2go library.
// This specific file contains all the code related to credentials.
package repository

import (
	"io/ioutil"
	"os"
	"os/user"
	"regexp"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/libgit2/git2go"
)

const (
	gitConfigFilename = ".gitconfig"
	gitConfigExpr     = `(?si)\[user\].+name\s=\s([\w\d ]+).+email\s=\s([\w\d@\.]+)`

	strategySSHAgent = "ssh-agent"
	strategySSHKey   = "ssh-key"
)

// Credentials contains the details of the user who's doing the push
// and the push strategy.
type Credentials struct {
	User *User
	Push *Push
}

// User represents the git user who will be used as signature
// for git operations.
type User struct {
	Name, Email string
}

// Push holds the configuration about the git push strategy.
type Push struct {
	Strategy, Username                string
	PublicKey, PrivateKey, Passphrase string
}

// SetCredentials sets the informations required for signing
// and pushing Git changes.
func (r *Repository) SetCredentials(creds *Credentials) {
	r.credentials = creds
}

// GetSignature returns the user infos required for
// signing Git changes.
func (r *Repository) GetSignature() *git.Signature {
	return &git.Signature{
		Name:  r.credentials.User.Name,
		Email: r.credentials.User.Email,
		When:  time.Now(),
	}
}

// credentialsCallback is linked to the git.RemoteCallbacks
func (r *Repository) credentialsCallback(url string, username string, allowedTypes git.CredType) (git.ErrorCode, *git.Cred) {
	var ret int
	var cred git.Cred

	switch r.credentials.Push.Strategy {
	case strategySSHKey:
		ret, cred = git.NewCredSshKey(r.credentials.Push.Username, r.credentials.Push.PublicKey, r.credentials.Push.PrivateKey, r.credentials.Push.Passphrase)
		break
	case strategySSHAgent:
		ret, cred = git.NewCredSshKeyFromAgent(r.credentials.Push.Username)
		break
	}

	return git.ErrorCode(ret), &cred
}

// certificateCheckCallback is linked to the git.RemoteCallbacks
func (r *Repository) certificateCheckCallback(cert *git.Certificate, valid bool, hostname string) git.ErrorCode {
	return 0
}

// GetUserFromGitConfig extracts from the local git config the
// username and email address of the git user if already configured.
func GetUserFromGitConfig() *User {
	usr, err := user.Current()
	if err != nil {
		log.Debug(err)
		return nil
	}

	gic := usr.HomeDir + string(os.PathSeparator) + gitConfigFilename
	if _, err := os.Stat(gic); os.IsNotExist(err) {
		log.Debug(err)
		return nil
	}

	dat, err := ioutil.ReadFile(gic)
	if err != nil {
		log.Debug(err)
		return nil
	}

	r, err := regexp.Compile(gitConfigExpr)
	if err != nil {
		log.Debug(err)
		return nil
	}

	res := r.FindStringSubmatch(string(dat))
	if len(res) == 0 {
		log.Debug(err)
		return nil
	}

	return &User{
		Name:  res[1],
		Email: res[2],
	}
}
