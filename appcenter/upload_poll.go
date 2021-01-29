package appcenter

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/pterm/pterm"
)

var path string

type uploadReleaseStatusResponse struct {
	ID                string `json:"id"`
	UploadStatus      string `json:"upload_status"`
	ErrorDetails      string `json:"error_details"`
	ReleaseDistinctID int64  `json:"release_distinct_id"`
}

// PollForRelease will poll AppCenter till the upload is ready to be pulished
func (s *UploadService) PollForRelease(ctx context.Context, uploadID string) (int64, error) {
	sp, err := pterm.DefaultSpinner.Start("Waiting for the release to be published")
	if err != nil {
		return -1, NewAppCenterError(PollingError, err)
	}

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
			sp.UpdateText(fmt.Sprintf("Waiting for the release to be published (Try: %d)", count))
			if count > 60 {
				sp.Fail()
				return -1, NewAppCenterError(PollingFailed, nil)
			}

			// polling for result
			c, err := s.poll(ctx)
			if err == nil {
				sp.Success(fmt.Sprintf("Release is ready to be published (ID: %d)", c))
				return c, nil
			}
		}
	}

	return -1, nil
}

// we are polling the release till it's status is "ready to be published"
func (s UploadService) poll(ctx context.Context) (int64, error) {
	var status uploadReleaseStatusResponse
	if err := s.client.NewAPIRequest(ctx, http.MethodGet, path, nil, &status); err != nil {
		return 0, err
	}

	// do we have the right upload status
	if status.UploadStatus == "readyToBePublished" {
		// if yes return the release distinct identifier
		return status.ReleaseDistinctID, nil
	}

	return 0, NewAppCenterError(PollingFailed, nil)
}
