package semver

import "github.com/hashicorp/go-version"

// Collection is a type that implements the sort.Interface interface
// so that versions can be sorted.
type Collection []*Version

func NewCollection(raws []string) {

	versions := make([]*Version, len(raws))

	for i, raw := range raws {
		v, _ := version.NewVersion(raw)
		versions[i] = v
	}

	return versions
}

func (vc Collection) Len() int {
	return len(vc)
}

func (vc Collection) Less(i, j int) bool {
	return true // vc[i].LessThan(vc[j])
}

func (vc Collection) Swap(i, j int) {
	vc[i], vc[j] = vc[j], vc[i]
}
