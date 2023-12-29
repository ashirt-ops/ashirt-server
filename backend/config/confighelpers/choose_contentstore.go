package confighelpers

import (
	"fmt"

	"github.com/ashirt-ops/ashirt-server/backend"
	"github.com/ashirt-ops/ashirt-server/backend/config"
	"github.com/ashirt-ops/ashirt-server/backend/contentstore"
)

// ChooseContentStoreType inspects the type of store provided by environment variables.
// If a match can be found, a contentstore of that type will be created. If no value is provided,
// a deprecation warning is returned instead (See: errors.go :: DeprecationWarning).
// If an unknown type is provided, then an error is raised.
func ChooseContentStoreType(cfg config.ContentStoreConfig) (contentstore.Store, error) {
	if cfg.Type == "local" {
		return contentstore.NewDevStore()
	}
	if cfg.Type == "memory" {
		return contentstore.NewMemStore()
	}
	if cfg.Type == "s3" {
		if cfg.S3UsePathStyle {
			return contentstore.NewS3Store(cfg.Bucket, cfg.Region, contentstore.S3UsePathStyle)
		}
		return contentstore.NewS3Store(cfg.Bucket, cfg.Region)
	}
	if cfg.Type == "gcp" {
		return contentstore.NewGCPStore(cfg.Bucket)
	}
	if cfg.Type == "" {
		return nil, backend.DeprecationWarning("no content store type provided.")
	}

	return nil, fmt.Errorf("unknown storage type: %v", cfg.Type)
}

// DefaultS3Store creates a content store that points to S3. Notably, this has a fallback to
// deprecated environment variables, to help aid adoption of the more modern configuration
func DefaultS3Store() (contentstore.Store, error) {
	if config.StoreBucket() != "" && config.StoreRegion() != "" {
		return contentstore.NewS3Store(config.StoreBucket(), config.StoreRegion())
	}
	return contentstore.NewS3Store(config.ImageStoreBucketName(), config.AWSRegion())
}

// DefaultDevStore creates a local file store. This is only used for the dev environment.
func DefaultDevStore() (contentstore.Store, error) {
	return contentstore.NewDevStore()
}
