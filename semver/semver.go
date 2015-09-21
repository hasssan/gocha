// Package semver provides some functions to help deal with semver version numbers.
// http://semver.org/
package semver

import (
	"fmt"
	"regexp"
	"strconv"
)

const (
	Expr   = `(?i)^v?(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(?:-[\da-z\-]+(\.[\da-z\-]+)*)?(\+[\da-z\-]+(?:\.[\da-z\-]+)*)?$`
	Format = "%s.%s.%s"
)

func IsValid(v string) bool {
	r, err := regexp.Compile(Expr)
	if err != nil {
		return false
	}
	res := r.MatchString(v)
	return res
}

func GetNextMajorVersion(v string) (string, error) {
	nv, err := getIncrementVersionByIndex(v, 1)
	if err != nil {
		return "", err
	}
	return nv, nil
}

func GetNextMinorVersion(v string) (string, error) {
	nv, err := getIncrementVersionByIndex(v, 2)
	if err != nil {
		return "", err
	}
	return nv, nil
}

func GetNextPatchVersion(v string) (string, error) {
	nv, err := getIncrementVersionByIndex(v, 3)
	if err != nil {
		return "", err
	}
	return nv, nil
}

func getIncrementVersionByIndex(v string, i int) (string, error) {
	var err error

	r, err := regexp.Compile(Expr)
	if err != nil {
		return "", err
	}

	res := r.FindStringSubmatch(v)

	vInt, err := strconv.Atoi(res[i])
	if err != nil {
		return "", err
	}
	res[i] = strconv.Itoa(vInt + 1)

	return fmt.Sprintf(Format, res[1], res[2], res[3]), nil
}
