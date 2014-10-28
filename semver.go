package semver

import (
	"fmt"
	"regexp"
	"strconv"
)

var (
	ErrorMalformedVersion = Error{"Malformed Version: %s", ""}
)

var (
	PreRegexp     *regexp.Regexp = regexp.MustCompile(`^([a-zA-Z1-9][0-9A-Za-z-]+)?(?:.([0-9]+))?`)
	productRegexp *regexp.Regexp = regexp.MustCompile(`^[a-zA-Z0-9-]+[a-zA-Z]$`)

	versionRegexp *regexp.Regexp = regexp.MustCompile(`^(?:([a-zA-Z0-9-]+[a-zA-Z])-)?` +
		`([0-9]+(?:\.[0-9]+){0,2})` +
		`(?:-([0-9A-Za-z]+(?:\.[0-9A-Za-z]+)*))?` +
		`(?:\+([0-9A-Za-z]+(?:\.[0-9A-Za-z]+)*))?$`)
)

type Pre struct {
	Ident string
	Patch int
}

func NewPre(raw string) *Pre {

	parts := PreRegexp.FindStringSubmatch(raw)

	i := parts[1]
	p, err := strconv.Atoi(parts[2])

	pre := Pre{i, p}

	if pre.Ident == "" || err != nil {
		return nil
	}
	return &pre
}

var stages = map[string]int{
	"pre":   0,
	"alpha": 1,
	"beta":  2,
	"rc":    3,
}

// Compares Pre v to o:
// -1 == v is less than o
// 0 == v is equal to o
// 1 == v is greater than o

func (v *Pre) Compare(o *Pre) int {
	if stages[v.Ident] > stages[o.Ident] {
		return 1
	}

	if stages[v.Ident] < stages[o.Ident] {
		return -1
	}

	if v.Patch > o.Patch {
		return 1
	}

	if v.Patch < o.Patch {
		return -1
	}

	return 0
}

type Version struct {
	Raw string

	Product string
	Major   uint16
	Minor   uint16
	Patch   uint16
	Pre     *Pre
	Meta    string
}

func NewVersion(raw string) (*Version, error) {

	parts := versionRegexp.FindStringSubmatch(raw)
	if parts == nil {
		return nil, ErrorMalformedVersion.Fault(raw)
	}

	ver := new(Version)

	ver.Raw = parts[0]

	ver.Product = parts[1]

	fmt.Sscanf(parts[2], "%d.%d.%d", &ver.Major, &ver.Minor, &ver.Patch)

	ver.Pre = NewPre(parts[3])

	ver.Meta = parts[4]

	return ver, nil
}

// Compares Versions v to o:
// -1 == v is less than o
// 0 == v is equal to o
// 1 == v is greater than o

func (v *Version) Compare(o *Version) int {

	//The MIT License

	if v.Major > o.Major {
		return 1
	}
	if v.Major < o.Major {
		return -1
	}

	if v.Minor > o.Minor {
		return 1
	}

	if v.Minor < o.Minor {
		return -1
	}

	if v.Patch > o.Patch {
		return 1
	}

	if v.Patch < o.Patch {
		return -1
	}

	return v.Pre.Compare(o.Pre)
}
