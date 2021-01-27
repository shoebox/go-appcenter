package appcenter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

type uploadReleaseBody struct {
	BuildVersion string `json:"build_version"`
	BuildNumber  string `json:"build_number"`
}

func (s *UploadService) SetMetaData(
	ctx context.Context,
	uploadDomain string,
	id string,
	fileName string,
	fileSize int64,
	token string,
	content_type string,
	buildNumber string,
	buildVersion string,
) (*metadataResponse, error) {
	log.Info().
		Str("Upload Domain", uploadDomain).
		Str("Token", token).
		Str("Content-Type", content_type).
		Str("File-Name", fileName).
		Msg("Applying meta data for upload")

	url := fmt.Sprintf(
		"%v/upload/set_metadata/%v?file_name=%v&file_size=%v&token=%v&content_type=%v",
		uploadDomain,
		id,
		fileName,
		fileSize,
		token,
		content_type,
	)

	// optional body
	body := uploadReleaseBody{}
	if buildNumber != "" {
		body.BuildNumber = buildNumber
	}

	if buildVersion != "" {
		body.BuildVersion = buildVersion
	}

	var m metadataResponse
	if _, err := s.client.simpleRequest(ctx, http.MethodPost, url, &body, &m); err != nil {
		return &m, NewAppCenterError(MetadataError, err)
	}

	log.Info().Msg("Metadata applied")

	return &m, nil
}
