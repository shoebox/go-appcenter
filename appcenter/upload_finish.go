package appcenter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pterm/pterm"
)

// FinishingUploadResponse response definition of the upload finished endpoint
type FinishingUploadResponse struct {
	Error        *bool   `json:"error,omitempty"`
	ChunkNum     *int64  `json:"chunk_num,omitempty"`
	ErrorCode    *string `json:"error_code,omitempty"`
	Message      *string `json:"message,omitempty"`
	Location     *string `json:"location,omitempty"`
	RawLocation  *string `json:"raw_location,omitempty"`
	AbsoluteURI  *string `json:"absolute_uri,omitempty"`
	State        *string `json:"state,omitempty"`
	UploadStatus *string `json:"uploadstatus,omitempty"`
}

// FinishingUpload will notify AppCenter that the upload is finished
func (s *UploadService) FinishingUpload(
	ctx context.Context,
	uploadDomain string,
	packageAssetID string,
	urlEncodedToken string,
	ID string,
) (*FinishingUploadResponse, error) {
	sp, err := pterm.DefaultSpinner.Start("Completing upload")
	if err != nil {
		return nil, err
	}

	var res FinishingUploadResponse

	url := fmt.Sprintf(
		"%v/upload/finished/%v?token=%v",
		uploadDomain,
		packageAssetID,
		urlEncodedToken,
	)

	_, err = s.client.simpleRequest(ctx, http.MethodPost, url, nil, &res)
	if err == nil {
		sp.Success()
	}

	return &res, err
}
