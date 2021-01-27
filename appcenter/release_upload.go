package appcenter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

type uploadReeleaseResponse struct {
	Status string `json:"upload_status"`
	ID     string `json:"upload_id"`
}

func (s *UploadService) ReleaseUpload(ctx context.Context, uploadID string) error {
	log.Info().
		Str("UploadID", uploadID).
		Msg("Releasing upload")

	var res uploadReeleaseResponse
	err := s.client.NewAPIRequest(
		ctx,
		http.MethodPatch,
		fmt.Sprintf("uploads/releases/%v", uploadID),
		commitUploadBody{Status: "uploadFinished", ID: uploadID},
		&res,
	)

	return err
}
