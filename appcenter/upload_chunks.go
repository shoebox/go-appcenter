package appcenter

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math"
	"net/http"
	"runtime"
	"sync"

	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type chunkUploadResponse struct {
	Error       bool   `json:"error,omitempty"`
	ChunkNum    int64  `json:"chunk_num,omitempty"`
	ErrorCode   string `json:"error_code,omitempty"`
	Message     string `json:"message,omitempty"`
	Location    string `json:"location,omitempty"`
	RawLocation string `json:"raw_location,omitempty"`
	AbsoluteURI string `json:"absolute_uri,omitempty"`
	State       string `json:"state,omitempty"`
}

type Chunk struct {
	ID   int
	URL  string
	Data []byte
}

// UploadChunks allow to upload a file by determined chunk size and count to AppCenter
func (s *UploadService) UploadChunks(
	ctx context.Context,
	reader io.Reader,
	uploadDomain string,
	uploadID string,
	packageAssetID string,
	urlEncodedToken string,
	chunkSize int,
	chunkCount int,
	fileSize int64,
	contentType string,
) error {
	var wg sync.WaitGroup

	jobc := make(chan Chunk, chunkCount)
	g, ctx := errgroup.WithContext(ctx)

	for i := 0; i < runtime.NumCPU(); i++ {
		go s.chunkUploadWorker(ctx, i, &wg, jobc, g)
	}

	for i := 0; i < chunkCount; i++ {
		c := Chunk{
			URL: fmt.Sprintf(
				"%s/upload/upload_chunk/%v?block_number=%v&token=%v",
				uploadDomain,
				packageAssetID,
				i+1,
				urlEncodedToken,
			),
			ID: (i + 1),
		}

		// chunk start/end position
		start := int64(i * chunkSize)
		end := start + int64(chunkSize)

		// module of the filesize
		if end > fileSize {
			end = int64(math.Min(float64(fileSize), float64(end)))
		}

		c.Data = make([]byte, end-start)

		// Reading the bytes
		if _, err := reader.Read(c.Data); err != nil {
			return NewAppCenterError(ChunkingError, err)
		}
		log.Info().Int("ChunkID", c.ID).Msg("Preparing chunk")

		// sending the chunk object to the worker pool
		jobc <- c
	}

	// closing the job channel
	close(jobc)

	wg.Wait()

	return g.Wait()
}

func (s *UploadService) chunkUploadWorker(
	ctx context.Context,
	id int,
	wg *sync.WaitGroup,
	jobs <-chan Chunk,
	g *errgroup.Group,
) {
	wg.Add(1)
	defer wg.Done()

	for j := range jobs {
		log.Info().Int("Chunk ID", j.ID).Msg("Chunk upload started")
		g.Go(func() error {

			r := chunkUploadResponse{}
			resp, err := s.client.simpleRequest(ctx, http.MethodPost, j.URL, bytes.NewBuffer(j.Data), &r)
			if err != nil {
				return err
			} else if resp.StatusError != nil {
				return resp.StatusError
			} else {
				log.Info().Int("Chunk ID", j.ID).Msg("Chunk upload complete")
			}

			return nil
		})
	}
}
