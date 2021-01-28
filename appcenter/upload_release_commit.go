package appcenter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pterm/pterm"
)

type commitUploadBody struct {
	Status    string `json:"upload_status"`
	ID        string `json:"upload_id"`
	ReleaseID int    `json:"release_distinct_id"`
}

type commitReleaseResponse struct {
	ID                string `json:"id"`
	UploadStatus      string `json:"upload_status"`
	ReleaseDistinctID int    `json:"release_distinct_id"`
}

func (s *UploadService) UploadCommitRelease(ctx context.Context, id string, uploadID string) (*string, error) {
	sp, err := pterm.DefaultSpinner.Start("Updating status of the release")
	if err != nil {
		return nil, err
	}

	var res commitReleaseResponse
	path := fmt.Sprintf("uploads/releases/%v", uploadID)

	//
	if err := s.client.NewAPIRequest(
		ctx,
		http.MethodPatch,
		path,
		commitUploadBody{Status: "uploadFinished", ID: uploadID},
		&res,
	); err != nil {
		sp.Fail()
		return nil, err
	}

	sp.Success(fmt.Sprintf("Release status update complete (Release ID: '%v')", res.ID))
	return &res.ID, nil
}
