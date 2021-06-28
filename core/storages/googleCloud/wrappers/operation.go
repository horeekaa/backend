package googlecloudstoragecorewrappers

import (
	"context"

	"cloud.google.com/go/storage"
	googlecloudstoragecorewrapperinterfaces "github.com/horeekaa/backend/core/storages/googleCloud/interfaces/wrappers"
)

type gcsObjectHandle struct {
	*storage.ObjectHandle
}

func (gcsHandle *gcsObjectHandle) NewWriter(ctx context.Context) googlecloudstoragecorewrapperinterfaces.GCSWriter {
	return gcsHandle.NewWriter(ctx)
}

func NewGCSObjectHandle(wrappedObjectHandle *storage.ObjectHandle) (googlecloudstoragecorewrapperinterfaces.GCSObjectHandle, error) {
	return &gcsObjectHandle{
		wrappedObjectHandle,
	}, nil
}
