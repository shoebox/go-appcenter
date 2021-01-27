package appcenter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

type commitUploadBody struct {
	Status    string `json:"upload_status"`
	ID        string `json:"upload_id"`
	ReleaseID int    `json:"release_distinct_id"`
}

type CommitReleaseResponse struct {
	ID                string `json:"id"`
	UploadStatus      string `json:"upload_status"`
	ReleaseDistinctId int    `json:"release_distinct_id"`
}

func (s *UploadService) UploadCommitRelease(ctx context.Context, id string, uploadID string) (*string, error) {
	log.Info().
		Str("ID", id).
		Str("UploadID", uploadID).
		Msg("Updating the status of the release upload")

	var res CommitReleaseResponse
	path := fmt.Sprintf("uploads/releases/%v", uploadID)

	//
	if err := s.client.NewAPIRequest(
		ctx,
		http.MethodPatch,
		path,
		commitUploadBody{Status: "uploadFinished", ID: uploadID},
		&res,
	); err != nil {
		return nil, err
	}
	return &res.ID, nil
}
