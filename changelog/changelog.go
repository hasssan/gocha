// Package changelog handles the changelog generation.
package changelog

import (
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"

	"github.com/flosch/pongo2"
	"github.com/jgautheron/gocha/message"
	"github.com/jgautheron/gocha/repository"
)

const (
	templateFile  = "template/changelog-template.md"
	changelogFile = "CHANGELOG.md"
)

// Generate will lookup the commits for the given tag and create a CHANGELOG.md file in the current path.
func Generate(rp *repository.Repository, tag string, appName string, outputFile string) {
	var err error
	var tg repository.Tag

	if tag != "" {
		tg, err = rp.GetTag(tag)
	} else {
		tg, err = rp.GetLastTag()
	}

	if err != nil {
		log.Fatal(err)
	}

	cmts, err := rp.GetCommitListForTag(tg)
	if err != nil {
		log.Fatal(err)
	}

	ms, err := message.GetMessageGroup(cmts)
	if err != nil {
		log.Fatal(err)
	}

	url, err := rp.GetOriginURL()
	if err != nil {
		log.Fatal(err)
	}

	output, err := getFilledTemplate(pongo2.Context{
		"appName":       appName,
		"version":       tg.Name,
		"message_group": ms,
		"url":           url,
	}, templateFile)
	if err != nil {
		log.Fatal(err)
	}

	fileInfo, err := os.Stat(outputFile)
	if err != nil {
		log.Fatal(err)
	}

	if fileInfo.IsDir() {
		if outputFile[len(outputFile)-1:] != "/" {
			outputFile += string(os.PathSeparator)
		}
		outputFile += changelogFile
	}

	err = ioutil.WriteFile(outputFile, output, 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("%s has been successfully created!", outputFile)
}

// getFilledTemplate returns the filled template as a slice of bytes.
// Initially wanted to use here the stdlib's text/template but ran into issues
// with the if instruction.
// The template looks quite ugly because of the blank lines left by the tags.
// https://code.djangoproject.com/ticket/2594 (WONTFIX)
// https://github.com/flosch/pongo2/issues/94
func getFilledTemplate(ctxt pongo2.Context, tplFile string) ([]byte, error) {
	t := pongo2.Must(pongo2.FromFile(tplFile))
	output, err := t.ExecuteBytes(ctxt)
	if err != nil {
		log.Fatal(err)
	}
	return output, nil
}
