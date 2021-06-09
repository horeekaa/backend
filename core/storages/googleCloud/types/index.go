package googlecloudstoragecoretypes

import (
	"io"
)

type GCSFileUpload struct {
	File        io.Reader
	Filename    string
	Size        int64
	ContentType string
}
