// Package bumper contains the logic for version bumping.
package bumper

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/jgautheron/codename-generator"
	"github.com/jgautheron/gocha/message"
	"github.com/jgautheron/gocha/repository"
	"github.com/jgautheron/gocha/semver"
)

const (
	Major = "major"
	Minor = "minor"
	Patch = "patch"
)

func Up(rp *repository.Repository, bmp string) {
	var err error

	lt, err := rp.GetLastTag()
	if err != nil {
		log.Fatal(err)
	}

	log.Debugf("Current tag is: %s", lt.Name)

	var nxt string
	switch string(bmp) {
	case Major:
		nxt, err = semver.GetNextMajorVersion(lt.Name)
		if err != nil {
			log.Fatal(err)
		}
		break
	case Minor:
		nxt, err = semver.GetNextMinorVersion(lt.Name)
		if err != nil {
			log.Fatal(err)
		}
		break
	case Patch:
		nxt, err = semver.GetNextPatchVersion(lt.Name)
		if err != nil {
			log.Fatal(err)
		}
		break
	}

	log.Debugf("Next tag is: %s", nxt)

	// Generate a codename
	codename, err := codename.Get(codename.Sanitized)
	if err != nil {
		log.Fatal(err)
	}
	log.Debugf("The generated codename is: %s", codename)

	msg, err := message.New(
		message.Chore,
		"release",
		fmt.Sprintf("v%s codename(%s)", nxt, codename),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = rp.CreateAndPushTag(nxt, msg.String())
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("The tag %s has been successfully pushed", nxt)
}
