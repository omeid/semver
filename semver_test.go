package semver

import (
	"testing"
)

func TestNewVersion(t *testing.T) {
	cases := []struct {
		version string
		err     bool
	}{
		{"package-ver-1.2.3", false},
		{"1.0", false},
		{"1", false},
		{"1.2.beta", true},
		{"foo", true},
		{"1.2-5", false},
		{"1.2-beta.5", false},
		{"\n1.2", true},
		{"1.2.0-x.Y.0+metadata", false},
		{"1.2.3.4", true},
	}

	for _, tc := range cases {
		v, err := NewVersion(tc.version)
		if tc.err && err == nil {
			t.Fatalf("expected error for version: %s \n %v", tc.version, v)
		} else if !tc.err && err != nil {
			t.Fatalf("error for version: %s \n %v", tc.version, err)
		}
	}
}
