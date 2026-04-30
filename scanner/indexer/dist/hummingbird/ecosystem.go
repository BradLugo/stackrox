package hummingbird

import (
	"context"

	"github.com/quay/claircore/indexer"
)

// NewEcosystem provides the DistributionScanner for the Hummingbird OS
// ecosystem. Package, repository, and matcher integration are intentionally
// omitted: this ecosystem is label-only.
func NewEcosystem(_ context.Context) *indexer.Ecosystem {
	return &indexer.Ecosystem{
		Name: "hummingbird",
		PackageScanners: func(_ context.Context) ([]indexer.PackageScanner, error) {
			return nil, nil
		},
		DistributionScanners: func(_ context.Context) ([]indexer.DistributionScanner, error) {
			return []indexer.DistributionScanner{new(DistributionScanner)}, nil
		},
		RepositoryScanners: func(_ context.Context) ([]indexer.RepositoryScanner, error) {
			return nil, nil
		},
	}
}