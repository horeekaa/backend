package googlecloudstoragecorewrapperinterfaces

import (
	"context"
	"io"
)

type GCSWriter interface {
	io.Writer
	Close() error
}

type GCSObjectHandle interface {
	NewWriter(ctx context.Context) GCSWriter
	Delete(ctx context.Context) error
}
