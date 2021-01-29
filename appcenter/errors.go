package appcenter

import "fmt"

const (
	// ChunkingError chunks upload step failed
	ChunkingError = "Failed to upload chunk"

	// InputFileError failed to validate the input file
	InputFileError = "Input file error"

	// MetadataError failed to apply metadata to the upload request
	MetadataError = "Apply metadata error"

	// PollingError failure while waiting for the upload to be ready to be published
	PollingError = "Timeout while waiting for upload to be ready to be published"

	// PollingFailed timeout while waiting for the upload to be ready to be published
	PollingFailed = "Polling failed"

	// UploadRequestError failed to request upload
	UploadRequestError = "Upload request error"
)

// AppCenterError generic error defintiion
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

// NewAppCenterError helper method to create a new AppCenterError
func NewAppCenterError(msg string, err error) error {
	return AppCenterError{msg: msg, err: err}
}
