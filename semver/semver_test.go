package semver_test

import (
	"testing"

	"github.com/jgautheron/gocha/semver"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSemverValidation(t *testing.T) {
	var v bool

	Convey("Valid versions should be validated", t, func() {
		v = semver.IsValid("1.0.0")
		So(v, ShouldBeTrue)

		v = semver.IsValid("1.0.0-foo")
		So(v, ShouldBeTrue)

		v = semver.IsValid("11.12.0-rc.1")
		So(v, ShouldBeTrue)
	})

	Convey("Invalid versions should be refused", t, func() {
		v = semver.IsValid("1.0")
		So(v, ShouldBeFalse)

		v = semver.IsValid("1")
		So(v, ShouldBeFalse)

		v = semver.IsValid("foo")
		So(v, ShouldBeFalse)
	})
}

func TestBumpMajorVersion(t *testing.T) {
	Convey("The major digit should be bumped", t, func() {
		mj, err := semver.GetNextMajorVersion("9.9.9")
		So(err, ShouldBeNil)
		So(mj, ShouldEqual, "10.9.9")
	})
}

func TestBumpMinorVersion(t *testing.T) {
	Convey("The major digit should be bumped", t, func() {
		mj, err := semver.GetNextMinorVersion("9.9.9")
		So(err, ShouldBeNil)
		So(mj, ShouldEqual, "9.10.9")
	})
}

func TestBumpPatchVersion(t *testing.T) {
	Convey("The major digit should be bumped", t, func() {
		mj, err := semver.GetNextPatchVersion("9.9.9")
		So(err, ShouldBeNil)
		So(mj, ShouldEqual, "9.9.10")
	})
}
