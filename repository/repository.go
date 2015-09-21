// Package repository wraps and simplifies the libgit2 bindings
// exposed in the git2go library.
package repository

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/jgautheron/gocha/semver"
	"github.com/libgit2/git2go"
)

var (
	errNoTagFound = errors.New("No semver tag has been found")
	errNoURLMatch = errors.New("No URL could be matched")
)

// Repository contains the original git.Repository object plus a few more
// useful things, such as the repository path on the FS, the credentials...
type Repository struct {
	path        string
	repository  *git.Repository
	credentials *Credentials
}

// Tag holds the information about a given tag.
type Tag struct {
	Name   string
	Date   time.Time
	Target *git.Oid
}

// Commit holds the information about a given commit.
type Commit struct {
	Description string
	Date        time.Time
	ID          *git.Oid
}

type timeSlice []Tag

// Forward request for length
func (p timeSlice) Len() int {
	return len(p)
}

// Define compare
func (p timeSlice) Less(i, j int) bool {
	return p[i].Date.Before(p[j].Date)
}

// Define swap over an array
func (p timeSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p timeSlice) GetSlice() []Tag {
	var ts []Tag
	for _, t := range p {
		ts = append(ts, t)
	}
	return ts
}

// New returns a new instance of Repository
func New(path string) (*Repository, error) {
	var err error

	// If the path is not defined, default to the current folder
	if len(path) == 0 {
		path, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}

	// Init the repo
	repository, err := git.OpenRepository(path)
	if err != nil {
		return nil, err
	}

	return &Repository{
		path:       path,
		repository: repository,
	}, nil
}

// GetRepository returns the Repository instance.
func (r *Repository) GetRepository() *git.Repository {
	return r.repository
}

// GetTags returns the tags list for the current repository.
func (r *Repository) GetTags() ([]Tag, error) {
	var ts []Tag

	// Retrieve the tags
	err := r.repository.Tags.Foreach(func(name string, id *git.Oid) error {
		tn := strings.Replace(name, "refs/tags/", "", -1)
		tg, err := r.buildTag(tn, id)
		if err != nil {
			return err
		}

		if semver.IsValid(tn) {
			ts = append(ts, tg)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	st := make(timeSlice, 0, len(ts))
	for _, d := range ts {
		st = append(st, d)
	}
	sort.Sort(st)

	if len(st) == 0 {
		return nil, errNoTagFound
	}

	return st.GetSlice(), nil
}

// GetLastTag returns the latest tag pushed.
func (r *Repository) GetLastTag() (Tag, error) {
	tags, err := r.GetTags()
	if err != nil {
		return Tag{}, err
	}

	return tags[len(tags)-1], nil
}

// CreateAndPushTag tags a repository and then pushes it automatically.
func (r *Repository) CreateAndPushTag(t string, msg string) error {
	var err error

	head, err := r.repository.Head()
	if err != nil {
		return err
	}
	defer head.Free()

	// Get latest commit id
	commit, err := r.repository.LookupCommit(head.Target())
	if err != nil {
		return err
	}
	defer commit.Free()

	if _, err = r.repository.Tags.Create(t, commit, r.GetSignature(), msg); err != nil {
		return err
	}

	// Retrieve the *Remote
	rm, err := r.repository.Remotes.Lookup("origin")
	if err != nil {
		return err
	}
	defer rm.Free()

	// Set the proper push URL
	// When using the credentials git.NewCredSshKey*, the URL must SSH typed:
	// user@repo.com/my/repo
	url, err := r.getSSHPushURL(rm.Url())
	if err != nil {
		return err
	}
	err = r.repository.Remotes.SetPushUrl("origin", url)
	if err != nil {
		return err
	}

	co := &git.PushOptions{
		RemoteCallbacks: git.RemoteCallbacks{
			CredentialsCallback:      r.credentialsCallback,
			CertificateCheckCallback: r.certificateCheckCallback,
		},
	}

	// Push the tag
	if err = rm.Push([]string{fmt.Sprintf("refs/tags/%s", t)}, co); err != nil {
		return err
	}

	return nil
}

// GetTag inspects the tag list and tries to match the given tag
// with an existing one.
func (r *Repository) GetTag(tag string) (Tag, error) {
	var err error
	var it *git.ReferenceIterator
	var tg Tag
	var emptyTag Tag

	if it, err = r.repository.NewReferenceIteratorGlob("refs/tags/*"); err != nil {
		return tg, err
	}
	defer it.Free()

	for {
		tr, err := it.Next()
		if err != nil {
			// The error here is empty
			break
		}

		// Match only the tag we're looking for
		tn := tr.Shorthand()
		if tn != tag {
			continue
		}

		tg, err = r.buildTag(tn, tr.Target())
		if err != nil {
			return emptyTag, err
		}

		break
	}

	if tg == emptyTag {
		return emptyTag, errNoTagFound
	}

	return tg, nil
}

// GetPreviousTagFor retrieves the previous Tag for the given Tag.
func (r *Repository) GetPreviousTagFor(tag Tag) (Tag, error) {
	var tg Tag

	tgs, err := r.GetTags()
	if err != nil {
		return tg, err
	}

	for idx, tg := range tgs {
		if tg.Name == tag.Name {
			// Cannot retrieve the previous tag if there is no previous tag
			if idx != 0 {
				return tgs[idx-1], nil
			}
		}
	}

	return tg, errNoTagFound
}

// GetCommitListForTag returns the list of commits associated
// with the given Tag.
func (r *Repository) GetCommitListForTag(tag Tag) ([]Commit, error) {
	var err error

	ptag, err := r.GetPreviousTagFor(tag)
	if err != nil {
		return nil, err
	}

	// Initialize and configure the rev walk
	rv, _ := r.repository.Walk()
	defer rv.Free()
	rv.Sorting(git.SortTime)

	// Start iterating from the tag's reference
	err = rv.Push(tag.Target)
	if err != nil {
		return nil, err
	}

	// Iterate until the previous tag
	err = rv.Hide(ptag.Target)
	if err != nil {
		return nil, err
	}

	var cmts []Commit

	var gi git.Oid
	for {
		err = rv.Next(&gi)
		if err != nil {
			// The error here is empty
			break
		}

		co, _ := r.repository.LookupCommit(&gi)
		cmts = append(cmts, Commit{
			Description: strings.TrimSpace(co.Message()),
			Date:        co.Committer().When,
			ID:          co.Id(),
		})
	}

	return cmts, nil
}

// buildTag creates a Tag from the given details.
func (r *Repository) buildTag(tn string, id *git.Oid) (Tag, error) {
	var cd time.Time

	// LookupTag will resolve only annotated tags
	tg, err := r.repository.LookupTag(id)
	if err != nil {
		// If LookupTag failed, then we're dealing with a
		// lightweight tag
		co, err := r.repository.LookupCommit(id)
		if err != nil {
			return Tag{}, err
		}

		cd = co.Committer().When
	} else {
		cd = tg.Tagger().When
	}

	return Tag{Name: tn, Date: cd, Target: id}, nil
}

// getSSHPushURL returns the given URL formatted for SSH.
func (r *Repository) getSSHPushURL(url string) (string, error) {
	if !strings.HasPrefix(url, "http") {
		return url, nil
	}

	rx, err := regexp.Compile(`^https?://([\w\-\.]+)/(.+)$`)
	if err != nil {
		return "", err
	}

	res := rx.FindStringSubmatch(url)
	if len(res) == 0 {
		return "", errNoURLMatch
	}

	return fmt.Sprintf("%s@%s:%s", r.credentials.Push.Username, res[1], res[2]), nil
}
