// Package hummingbird provides a ClairCore DistributionScanner that detects
// Hummingbird OS layers. It is label-only: no matcher or updater is registered,
// so OS-level CVEs are not produced for Hummingbird at this time.
package hummingbird

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"regexp"
	"runtime/trace"

	"github.com/quay/claircore"
	"github.com/quay/claircore/indexer"
	"github.com/quay/zlog"
)

var (
	_ indexer.DistributionScanner = (*DistributionScanner)(nil)
	_ indexer.VersionedScanner    = (*DistributionScanner)(nil)

	idLine        = regexp.MustCompile(`(?m)^ID="?hummingbird"?\s*$`)
	versionIDLine = regexp.MustCompile(`(?m)^VERSION_ID="?([^"\n]+?)"?\s*$`)
)

const osReleasePath = `etc/os-release`

// DistributionScanner detects Hummingbird OS by reading etc/os-release and
// matching ID=hummingbird. The version is taken verbatim from VERSION_ID.
//
// DistributionScanner can be used concurrently.
type DistributionScanner struct{}

// Name implements [indexer.VersionedScanner].
func (*DistributionScanner) Name() string { return "hummingbird" }

// Version implements [indexer.VersionedScanner]. Bump when Scan output for a
// given input changes, otherwise existing manifests will not be re-scanned.
func (*DistributionScanner) Version() string { return "1" }

// Kind implements [indexer.VersionedScanner].
func (*DistributionScanner) Kind() string { return "distribution" }

// Scan implements [indexer.DistributionScanner].
func (ds *DistributionScanner) Scan(ctx context.Context, l *claircore.Layer) ([]*claircore.Distribution, error) {
	defer trace.StartRegion(ctx, "Scanner.Scan").End()
	ctx = zlog.ContextWithValues(ctx,
		"component", "hummingbird/DistributionScanner.Scan",
		"version", ds.Version(),
		"layer", l.Hash.String())
	zlog.Debug(ctx).Msg("start")
	defer zlog.Debug(ctx).Msg("done")

	sys, err := l.FS()
	if err != nil {
		return nil, fmt.Errorf("hummingbird: unable to open layer: %w", err)
	}
	d, err := findDistribution(sys)
	if err != nil {
		return nil, fmt.Errorf("hummingbird: unexpected error reading files: %w", err)
	}
	if d == nil {
		zlog.Debug(ctx).Msg("layer is not Hummingbird OS")
		return nil, nil
	}
	return []*claircore.Distribution{d}, nil
}

func findDistribution(sys fs.FS) (*claircore.Distribution, error) {
	b, err := fs.ReadFile(sys, osReleasePath)
	switch {
	case errors.Is(err, nil):
	case errors.Is(err, fs.ErrNotExist):
		return nil, nil
	default:
		return nil, err
	}
	if !idLine.Match(b) {
		return nil, nil
	}
	m := versionIDLine.FindSubmatch(b)
	if m == nil || len(m[1]) == 0 {
		return nil, nil
	}
	return mkRelease(string(m[1])), nil
}