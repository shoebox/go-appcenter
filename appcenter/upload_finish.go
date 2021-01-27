package appcenter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

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

func (s *UploadService) FinishingUpload(
	ctx context.Context,
	uploadDomain string,
	packageAssetID string,
	urlEncodedToken string,
	ID string,
) (*FinishingUploadResponse, error) {
	log.Info().Msg("Finishing upload")

	var res FinishingUploadResponse

	url := fmt.Sprintf(
		"%v/upload/finished/%v?token=%v",
		uploadDomain,
		packageAssetID,
		urlEncodedToken,
	)

	_, err := s.client.simpleRequest(ctx, http.MethodPost, url, nil, &res)
	if err != nil {
		return &res, err
	}

	return &res, err
}
