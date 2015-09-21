package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringMessage(t *testing.T) {
	var m *Message
	var err error

	assert := assert.New(t)

	var strTests = []struct {
		in       []interface{}
		expected string
	}{
		{[]interface{}{Feat, "main", "simplify isolate scope bindings"}, "feat(main): simplify isolate scope bindings"},
		{[]interface{}{Fix, "semver", "fix the alpha/beta/... matching"}, "fix(semver): fix the alpha/beta/... matching"},
		{[]interface{}{Test, "foo", "add mock for the HTTP library"}, "test(foo): add mock for the HTTP library"},
		{[]interface{}{Refactor, "main", "export the business logic from the main package"}, "refactor(main): export the business logic from the main package"},
		{[]interface{}{Docs, "foo", "add readme for documenting the use cases"}, "docs(foo): add readme for documenting the use cases"},
		{[]interface{}{Style, "main", "format the code with gofmt"}, "style(main): format the code with gofmt"},
		{[]interface{}{Chore, "", "simplify isolate scope bindings"}, "chore: simplify isolate scope bindings"},
	}

	for _, tt := range strTests {
		m, err = New(tt.in[0].(messageType), tt.in[1].(string), tt.in[2].(string))
		assert.Nil(err)
		assert.IsType(&Message{}, m)
		assert.Equal(m.String(), tt.expected)
	}
}

func TestFullMessage(t *testing.T) {
	assert := assert.New(t)

	fulldesc := `feat($compile): simplify isolate scope bindings

Changed the isolate scope binding options to:
  - @attr - attribute binding (including interpolation)
  - =model - by-directional model binding
  - &expr - expression execution binding

This change simplifies the terminology as well as
number of choices available to the developer. It
also supports local name aliasing from the parent.

BREAKING CHANGE: isolate scope bindings definition has changed and
the inject option for the directive controller injection was removed.`
	msg, err := getMessageFromString(fulldesc)

	assert.Nil(err)
	assert.NotEmpty(msg.Type)
	assert.NotEmpty(msg.Scope)
	assert.NotEmpty(msg.Subject)
	assert.NotEmpty(msg.Body)
}

func TestShortMessage(t *testing.T) {
	assert := assert.New(t)

	shortdesc := `chore: update the readme to include notes about installation on OSX machines`
	msg, err := getMessageFromString(shortdesc)

	assert.Nil(err)
	assert.NotEmpty(msg.Type)
	assert.NotEmpty(msg.Subject)
	assert.Empty(msg.Scope)
	assert.Empty(msg.Body)
}
