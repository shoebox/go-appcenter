package appcenter

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

var path string

type UploadReleaseStatusResponse struct {
	ID                string `json:"id"`
	UploadStatus      string `json:"upload_status"`
	ErrorDetails      string `json:"error_details"`
	ReleaseDistinctId int64  `json:"release_distinct_id"`
}

// PollForRelease will poll AppCenter till the upload is ready to be pulished
func (s *UploadService) PollForRelease(ctx context.Context, uploadID string) (int64, error) {
	log.Info().Msg("Polling for release being ready to be published")

	// the path to poll against
	path = fmt.Sprintf("uploads/releases/%v", uploadID)

	// timer ticket every second
	t := time.NewTicker(time.Second)
	defer t.Stop()

	count := 0
	for range time.Tick(2 * time.Second) {
		select {
		// context cancellation handling
		case <-ctx.Done():
			return -1, NewAppCenterError(PollingError, nil)
		default:
			count++

			if count > 60 {
				return -1, NewAppCenterError(PollingFailed, nil)
			} else {
				// polling for result
				c, err := s.poll(ctx)
				if err == nil {
					return c, nil
				}
			}
		}
	}

	return -1, nil
}

// we are polling the release till it's status is "ready to be published"
func (s UploadService) poll(ctx context.Context) (int64, error) {
	log.Info().Msg("Polling for status change")
	var status UploadReleaseStatusResponse
	if err := s.client.NewAPIRequest(ctx, http.MethodGet, path, nil, &status); err != nil {
		return 0, err
	}

	// do we have the right upload status
	if status.UploadStatus == "readyToBePublished" {
		// if yes return the release distinct identifier
		log.Info().Msg("Release ready to be published")
		return status.ReleaseDistinctId, nil
	} else {
		log.Debug().Msg("Release not yet ready")
	}

	return 0, NewAppCenterError(PollingFailed, nil)
}
