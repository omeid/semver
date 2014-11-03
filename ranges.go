package semver

import (
	"fmt"
	"regexp"
	"strings"
)

// Range represents a single range for a version, such as
// ">= 1.0".
type Range struct {
	f        rangeFunc
	check    *Version
	original string
}

// Ranges is a slice of ranges. We make a custom type so that
// we can add methods to it.
type Ranges []*Range
type rangeFunc func(v, c *Version) bool

var rangeOperators map[string]rangeFunc
var rangeRegexp *regexp.Regexp

func init() {
	rangeOperators = map[string]rangeFunc{
		"":   rangeEqual,
		"=":  rangeEqual,
		"!=": rangeNotEqual,
		">":  rangeGreaterThan,
		"<":  rangeLessThan,
		">=": rangeGreaterThanEqual,
		"<=": rangeLessThanEqual,
		"*":  rangeAny,
	}

	ops := make([]string, 0, len(rangeOperators))
	for k, _ := range rangeOperators {
		ops = append(ops, regexp.QuoteMeta(k))
	}

	rangeRegexp = regexp.MustCompile(fmt.Sprintf(
		`(%s) *(%s)$`,
		strings.Join(ops, "|"),
		versionRegexpRaw))

}

// NewRange will parse one or more ranges from the given
// range string. The string must be a comma-separated list of
// ranges.
func NewRange(v string) (Ranges, error) {

	//Make a valid version.
	if v == "*" {
		v = "* 0.0.0"
	}

	vs := strings.Split(v, ",")
	result := make([]*Range, len(vs))
	for i, single := range vs {
		c, err := parseSingle(single)
		if err != nil {
			return nil, err
		}
		result[i] = c
	}
	return Ranges(result), nil
}

// Check tests if a version satisfies all the ranges.
func (cs Ranges) Check(v *Version) bool {
	for _, c := range cs {
		if !c.Check(v) {
			return false
		}
	}
	return true
}

// Returns the string format of the ranges
func (cs Ranges) String() string {
	csStr := make([]string, len(cs))
	for i, c := range cs {
		csStr[i] = c.String()
	}
	return strings.Join(csStr, ",")
}

// Check tests if a range is validated by the given version.
func (c *Range) Check(v *Version) bool {
	return c.f(v, c.check)
}
func (c *Range) String() string {
	return c.original
}
func parseSingle(v string) (*Range, error) {
	matches := rangeRegexp.FindStringSubmatch(v)
	if matches == nil {
		return nil, fmt.Errorf("Malformed range: %s", v)
	}
	check, err := NewVersion(matches[2])
	if err != nil {
		// This is a panic because the regular expression above should
		// properly validate any version.
		panic(err)
	}
	return &Range{
		f:        rangeOperators[matches[1]],
		check:    check,
		original: v,
	}, nil
}

//-------------------------------------------------------------------
// Range functions
//-------------------------------------------------------------------
func rangeAny(v, c *Version) bool {
	return true
}

func rangeEqual(v, c *Version) bool {
	return v.Equal(c)
}
func rangeNotEqual(v, c *Version) bool {
	return !v.Equal(c)
}
func rangeGreaterThan(v, c *Version) bool {
	return v.Compare(c) == 1
}
func rangeLessThan(v, c *Version) bool {
	return v.Compare(c) == -1
}
func rangeGreaterThanEqual(v, c *Version) bool {
	return v.Compare(c) >= 0
}
func rangeLessThanEqual(v, c *Version) bool {
	return v.Compare(c) <= 0
}
