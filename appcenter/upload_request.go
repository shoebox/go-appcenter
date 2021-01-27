package appcenter

import (
	"fmt"
	"os"
	"path/filepath"
)

// UploadTask wrap the required arguments for the upload specifications
type UploadTask struct {
	//APIToken  string
	AppName    string
	OwnerName  string
	FilePath   string
	Distribute DistributionPayload
	Option     ReleaseUploadPayload
}

func (r UploadTask) validateRequest() error {
	// validation of the sourece file to upload
	if err := r.validateSource(); err != nil {
		return err
	}

	// validation of the request settings
	return r.validateRequestBuildVersion()
}

func (r UploadTask) validateSource() error {
	_, err := os.Stat(r.FilePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("File `%v` does not exsts", r.FilePath)
	}
	return err
}

func (r UploadTask) validateRequestBuildVersion() error {
	ext := filepath.Ext(r.FilePath)

	if r.Option.BuildNumber == "" || r.Option.BuildVersion == "" {
		if ext == ".pkg" || ext == ".dmg" {
			return fmt.Errorf("'--build_version' and '--build_number' parameters "+
				"must be specified to upload file of extension %v", ext)
		}
	}

	if r.Option.BuildVersion == "" {
		if ext == ".zip" || ext == ".msi" {
			return fmt.Errorf("'--build_version' parameter must be "+
				"specified to upload fle of extension '%v'", ext)
		}
	}

	return nil
}
