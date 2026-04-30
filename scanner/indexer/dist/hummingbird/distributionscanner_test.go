package hummingbird

import (
	"io/fs"
	"os"
	"path"
	"strings"
	"testing"
)

func TestFindDistribution(t *testing.T) {
	root := os.DirFS("testdata/releasefiles")
	ents, err := fs.ReadDir(root, ".")
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range ents {
		t.Run(e.Name(), func(t *testing.T) {
			n := path.Base(t.Name())
			sub, err := fs.Sub(root, n)
			if err != nil {
				t.Fatal(err)
			}
			d, err := findDistribution(sub)
			if err != nil {
				t.Fatal(err)
			}
			switch {
			case strings.HasPrefix(n, "not-"):
				if d != nil {
					t.Fatalf("unexpected distribution: %s:%s", d.DID, d.VersionID)
				}
			default:
				if d == nil {
					t.Fatal("missing distribution")
				}
				if got, want := d.DID, "hummingbird"; got != want {
					t.Errorf("DID got %q, want %q", got, want)
				}
				if got, want := d.VersionID, n; got != want {
					t.Errorf("VersionID got %q, want %q", got, want)
				}
			}
		})
	}
}

func TestFindDistributionMissingOSRelease(t *testing.T) {
	d, err := findDistribution(fs.FS(emptyFS{}))
	if err != nil {
		t.Fatal(err)
	}
	if d != nil {
		t.Fatalf("unexpected distribution: %+v", d)
	}
}

type emptyFS struct{}

func (emptyFS) Open(string) (fs.File, error) { return nil, fs.ErrNotExist }

func TestMkReleaseMemoizes(t *testing.T) {
	a := mkRelease("20251124")
	b := mkRelease("20251124")
	if a != b {
		t.Fatalf("expected pointer equality for same VersionID")
	}
	c := mkRelease("20260101")
	if a == c {
		t.Fatalf("expected different pointers for different VersionIDs")
	}
}