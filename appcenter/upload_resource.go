package appcenter

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
)

type UploadResourceResponse struct {
	ID              string `json:"id,omitempty"`
	PackageAssetID  string `json:"package_asset_id,omitempty"`
	UploadDomain    string `json:"upload_domain,omitempty"`
	Token           string `json:"token,omitempty"`
	URLEncodedToken string `json:"url_encoded_token,omitempty"`
}

func (s *UploadService) RequestUploadResource(ctx context.Context, r UploadTask) (*UploadResourceResponse, error) {
	log.Info().Msg("Requesting upload resource")

	var result UploadResourceResponse

	if err := s.client.NewAPIRequest(
		ctx,
		http.MethodPost,
		"uploads/releases",
		r.Option,
		&result,
	); err != nil {
		return nil, NewAppCenterError(UploadRequestError, err)
	}

	log.Info().
		Str("Domain", result.UploadDomain).
		Str("ID", result.ID).
		Msg("Upload requested successfully")

	return &result, nil
}
