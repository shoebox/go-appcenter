package appcenter

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/pterm/pterm"
)

type releaseInfoResponseBody struct {
	ID                            int64  `json:"id,omitempty"`
	AppName                       string `json:"app_name,omitempty"`
	AppDisplayName                string `json:"app_display_name,omitempty"`
	AppOS                         string `json:"app_os,omitempty"`
	Version                       string `json:"version,omitempty"`
	Origin                        string `json:"origin,omitempty"`
	ShortVersion                  string `json:"short_version,omitempty"`
	ReleaseNotes                  string `json:"release_notes,omitempty"`
	ProvisioningProfileName       string `json:"provisioning_profile_name,omitempty"`
	ProvisioningProfileType       string `json:"provisioning_profile_type,omitempty"`
	ProvisioningProfileExpiryDate string `json:"provisioning_profile_expiry_date,omitempty"`
	IsProvisioningProfileSyncing  bool   `json:"is_provisioning_profile_syncing,omitempty"`
	Size                          int64  `json:"size,omitempty"`
	MinOS                         string `json:"min_os,omitempty"`
	DeviceFamily                  string `json:"device_family,omitempty"`
	AndroidMinAPILevel            string `json:"android_min_api_level,omitempty"`
	BundleIdentifier              string `json:"bundle_identifier,omitempty"`
	Fingerprint                   string `json:"fingerprint,omitempty"`
	UploadedAt                    string `json:"uploaded_at,omitempty"`
	DownloadURL                   string `json:"download_url,omitempty"`
	AppIconURL                    string `json:"app_icon_url,omitempty"`
	InstallURL                    string `json:"install_url,omitempty"`
	DestinationType               string `json:"destination_type,omitempty"`
	IsUdidProvisioned             bool   `json:"is_udid_provisioned,omitempty"`
	CanResign                     bool   `json:"can_resign,omitempty"`
	Enabled                       bool   `json:"enabled,omitempty"`
	Status                        string `json:"status,omitempty"`
	IsExternalBuild               bool   `json:"is_external_build,omitempty"`
}

func (s *UploadService) UploadResult(ctx context.Context, id int64) error {
	sp, err := pterm.DefaultSpinner.Start("Requesting the release details")
	if err != nil {
		return err
	}

	var res releaseInfoResponseBody
	path := fmt.Sprintf("releases/%v", id)

	//
	if err := s.client.NewAPIRequest(
		ctx,
		http.MethodGet,
		path,
		nil,
		&res,
	); err != nil {
		sp.Fail()
		return err
	}

	sp.Success()

	fields := reflect.TypeOf(res)
	values := reflect.ValueOf(res)

	num := fields.NumField()

	data := [][]string{}

	for i := 0; i < num; i++ {
		field := fields.Field(i)
		value := values.Field(i)

		key := field.Name
		val := ""

		switch value.Kind() {
		case reflect.String:
			val = value.String()
		case reflect.Bool:
			if value.Bool() {
				val = "YES"
			} else {
				val = "NO"
			}

		case reflect.Int64:
			val = strconv.FormatInt(value.Int(), 10)
		default:
		}

		if val != "" && len(val) < 60 {
			data = append(data, []string{key, val})
		}
	}

	pterm.DefaultTable.WithData(data).Render()

	return nil
}
