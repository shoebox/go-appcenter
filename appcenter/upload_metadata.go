package appcenter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pterm/pterm"
)

type uploadReleaseBody struct {
	BuildVersion string `json:"build_version"`
	BuildNumber  string `json:"build_number"`
}

// SetMetaData will apply meta data to the upload slot
func (s *UploadService) SetMetaData(
	ctx context.Context,
	uploadDomain string,
	id string,
	fileName string,
	fileSize int64,
	token string,
	contentType string,
	buildNumber string,
	buildVersion string,
) (*MetadataResponse, error) {
	sp, err := pterm.DefaultSpinner.Start("Applying meta-data")
	if err != nil {
		return nil, NewAppCenterError(UploadRequestError, nil)
	}

	url := fmt.Sprintf(
		"%v/upload/set_metadata/%v?file_name=%v&file_size=%v&token=%v&content_type=%v",
		uploadDomain,
		id,
		fileName,
		fileSize,
		token,
		contentType,
	)

	// optional body
	body := uploadReleaseBody{}
	if buildNumber != "" {
		body.BuildNumber = buildNumber
	}

	if buildVersion != "" {
		body.BuildVersion = buildVersion
	}

	var m MetadataResponse
	if _, err := s.client.simpleRequest(ctx, http.MethodPost, url, &body, &m); err != nil {
		return &m, NewAppCenterError(MetadataError, err)
	}

	sp.UpdateText("Metadata applied succesfully")
	sp.Success()

	return &m, nil
}
