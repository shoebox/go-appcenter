package appcenter

import (
	"context"
	"os"
	"path/filepath"
)

// UploadService definition
type UploadService struct {
	client *Client
}

type DistributionPayload struct {
	GroupName string
}

// ReleaseUploadPayload wrap optional informations about the release
type ReleaseUploadPayload struct {
	ReleaseID    int    `json:"release_id,omitempty"`
	BuildVersion string `json:"build_version,omitempty"`
	BuildNumber  string `json:"build_number,omitempty"`
}

type metadataResponse struct {
	Error          *bool   `json:"error,omitempty"`
	ID             *string `json:"id,omitempty"`
	ChunkSize      *int    `json:"chunk_size,omitempty"`
	ResumeRestart  *bool   `json:"resume_restart,omitempty"`
	ChunkList      []int64 `json:"chunk_list,omitempty"`
	BlobPartitions *int64  `json:"blob_partitions,omitempty"`
	StatusCode     *string `json:"status_code,omitempty"`
}

// Do start the upload request witht the provided parameters
func (s *UploadService) Do(ctx context.Context, r UploadTask) (int64, error) {
	if err := r.validateRequest(); err != nil {
		return -1, err
	}

	// Request Upload "slot"
	ur, err := s.RequestUploadResource(ctx, r)
	if err != nil {
		return -1, err
	}

	// convert to absolute path
	p, err := filepath.Abs(r.FilePath)
	if err != nil {
		return -1, NewAppCenterError(InputFileError, err)
	}
	content_type := ResolveContentType(filepath.Ext(p))

	// get target file infos
	fi, err := os.Stat(p)
	if err != nil {
		return -1, NewAppCenterError(InputFileError, err)
	}

	// Metadatas
	meta, err := s.SetMetaData(
		ctx,
		ur.UploadDomain,
		ur.PackageAssetID,
		fi.Name(),
		fi.Size(),
		ur.URLEncodedToken,
		content_type,
		r.Option.BuildNumber,
		r.Option.BuildVersion,
	)
	if err != nil {
		return -1, err
	}

	// Opening the file
	reader, err := os.Open(r.FilePath)
	if err != nil {
		return -1, err
	}
	defer reader.Close()

	// Uploading chunks
	err = s.UploadChunks(
		ctx,
		reader,
		ur.UploadDomain,
		ur.ID,
		ur.PackageAssetID,
		ur.URLEncodedToken,
		*meta.ChunkSize,
		len(meta.ChunkList),
		fi.Size(),
		content_type,
	)

	if err != nil {
		return -1, err
	}

	// finishing upload
	_, err = s.FinishingUpload(ctx, ur.UploadDomain, ur.PackageAssetID, ur.URLEncodedToken, ur.ID)
	if err != nil {
		return -1, err
	}

	// Committing release
	_, err = s.UploadCommitRelease(ctx, *meta.ID, ur.ID)
	if err != nil {
		return -1, err
	}

	rdid, err := s.PollForRelease(ctx, ur.ID)
	if err != nil {
		return -1, err
	}

	return rdid, nil
}

type CommitUploadBody struct {
	Status string `json:"upload_status"`
	ID     string `json:"upload_id"`
}

type CommitUploadResponse struct {
	ID     string `json:"id"`
	Status string `json:"upload_status"`
}
