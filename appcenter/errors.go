package appcenter

import "fmt"

const (
	ChunkingError      = "Failed to upload chunk"
	InputFileError     = "Input file error"
	MetadataError      = "Apply metadata error"
	PollingError       = "Timeout while waiting for upload to be ready to be published"
	PollingFailed      = "Polling failed"
	UploadRequestError = "Upload request error"
)

type AppCenterError struct {
	msg string
	err error
}

func (k AppCenterError) Error() string {
	if k.err != nil {
		return fmt.Sprintf("AppCenter error: %v (%v)", k.msg, k.err)
	}

	return k.msg
}

func NewAppCenterError(msg string, err error) error {
	return AppCenterError{msg: msg, err: err}
}
