package appcenter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pterm/pterm"
)

// UploadResourceResponse response body
type UploadResourceResponse struct {
	ID              string `json:"id,omitempty"`
	PackageAssetID  string `json:"package_asset_id,omitempty"`
	UploadDomain    string `json:"upload_domain,omitempty"`
	Token           string `json:"token,omitempty"`
	URLEncodedToken string `json:"url_encoded_token,omitempty"`
}

// RequestUploadResource will request appcenter for a new resouce assignement ready for upload
func (s *UploadService) RequestUploadResource(ctx context.Context, r UploadTask) (*UploadResourceResponse, error) {
	sp, err := pterm.DefaultSpinner.Start("Requesting upload ressource")
	if err != nil {
		return nil, NewAppCenterError(UploadRequestError, nil)
	}

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

	sp.UpdateText(fmt.Sprintf("Upload requested successfully (ID : %v", result.ID))
	sp.Success()

	return &result, nil
}
