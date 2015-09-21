// Package message wraps the commit message style logic.
//
// For now there is only one convention available, the AngularJS one.
// Read more about it here:
// https://docs.google.com/document/d/1QrDFcIiPjSLDn3EL15IJygNPiHORgU1_OOAqWjiDU5Y/edit
package message

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/jgautheron/gocha/repository"
)

var (
	errNoMatch     = errors.New("Could not match the message")
	errInvalidType = errors.New("The given message type was not recognised")
)

const (
	simpleFormat   = "%s: %s"
	extendedFormat = "%s(%s): %s"
	formatExpr     = `(?si)^([a-z]{3,})(?:\(([\w\d\-_$]+)\))?:([\w\d$@\:\(\)\-\.,'"=&_/\\ ]+)(.+)?`

	// Available Types
	NA messageType = iota
	Chore
	Test
	Docs
	Feat
	Fix
	Refactor
	Style
)

type MessageGroup struct {
	Type, Scope string
}

type messageType byte

func (m messageType) String() string {
	tp, _ := map[messageType]string{
		Chore:    "chore",
		Test:     "test",
		Docs:     "docs",
		Feat:     "feat",
		Refactor: "refactor",
		Style:    "style",
		Fix:      "fix",
	}[m]

	return tp
}

type Message struct {
	Type    messageType
	Scope   string
	Subject string
	Body    string

	Date time.Time
	ID   string
}

func New(tp interface{}, scope string, subj string) (*Message, error) {
	var err error

	switch tp.(type) {
	case string:
		tp, err = parseMessageType(tp.(string))
		if err != nil {
			return nil, err
		}
	}

	return &Message{
		Type:    tp.(messageType),
		Scope:   scope,
		Subject: subj,
	}, nil
}

// String returns a properly formatted commit message.
func (m *Message) String() string {
	tp := m.Type.String()

	if len(m.Scope) == 0 {
		return fmt.Sprintf(simpleFormat, tp, m.Subject)
	}

	return fmt.Sprintf(extendedFormat, tp, m.Scope, m.Subject)
}

// GetMessageGroup analyses the given commits, and returns a list of
// properly deconstructed messages for the current convention.
func GetMessageGroup(cmts []repository.Commit) (map[string]map[string][]Message, error) {
	ms := make(map[string]map[string][]Message)

	for _, co := range cmts {
		msg, err := getMessageFromString(co.Description)
		if err != nil {
			continue
		}

		msg.ID = co.ID.String()
		msg.Date = co.Date

		tstr := msg.Type.String()

		if len(msg.Scope) == 0 {
			msg.Scope = "none"
		}

		// Init the map
		_, ok := ms[tstr]
		if !ok {
			ms[tstr] = make(map[string][]Message)
		}

		ms[tstr][msg.Scope] = append(ms[tstr][msg.Scope], *msg)
	}

	return ms, nil
}

// getMessageFromString analyses the given commit message,
// and creates a Message out of it.
func getMessageFromString(msg string) (*Message, error) {
	// Initally wanted to use fmt.Sscanf to split the string
	// but it didn't work as expected: scanning requires items
	// to be space-separated
	r, err := regexp.Compile(formatExpr)
	if err != nil {
		return nil, err
	}

	res := r.FindStringSubmatch(msg)
	if len(res) != 5 {
		return nil, errNoMatch
	}

	tp := res[1]
	scope := res[2]
	subj := strings.TrimSpace(res[3])
	body := strings.TrimSpace(res[4])

	stp, err := parseMessageType(tp)
	if err != nil {
		return nil, err
	}

	return &Message{
		Type:    stp,
		Scope:   scope,
		Subject: subj,
		Body:    body,
	}, nil
}

// parseMessageType
func parseMessageType(msg string) (messageType, error) {
	// Validate the incoming message type
	tp, ok := map[string]messageType{
		"chore":    Chore,
		"test":     Test,
		"docs":     Docs,
		"feat":     Feat,
		"refactor": Refactor,
		"style":    Style,
		"fix":      Fix,
	}[msg]

	if !ok {
		return NA, errInvalidType
	}

	return tp, nil
}
