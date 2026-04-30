package hummingbird

import (
	"sync"

	"github.com/quay/claircore"
)

// relMap memoizes the Distributions handed out by this package so that
// pointer-equality continues to work across calls for the same VERSION_ID.
var relMap sync.Map

func mkRelease(versionID string) *claircore.Distribution {
	if v, ok := relMap.Load(versionID); ok {
		return v.(*claircore.Distribution)
	}
	v, _ := relMap.LoadOrStore(versionID, &claircore.Distribution{
		Name:       "Hummingbird OS",
		Version:    versionID,
		VersionID:  versionID,
		DID:        "hummingbird",
		PrettyName: "Hummingbird OS " + versionID,
	})
	return v.(*claircore.Distribution)
}