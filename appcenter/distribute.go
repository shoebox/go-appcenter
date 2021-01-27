package appcenter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

// DistributeService definition
type DistributeService struct {
	client *Client
}

type distributionGroupResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Origin string `json:"origin"`
}

type distributionBody struct {
	ID              string `json:"id"`
	MandatoryUpdate bool   `json:"mandatory_update"`
	NotifyTester    bool   `json:"notify_testers"`
}

type distributionResponse struct {
	ID                    string `json:"id"`
	MandatoryUpdate       bool   `json:"mandatory_update"`
	ProvisioningStatusURL string `json:"provisioning_status_url"`
}

// Do Distribute the designated release into the provided configuration
func (s *DistributeService) Do(ctx context.Context, releaseID int64, request UploadTask) error {
	if request.Distribute.GroupName != "" {
		log.Info().Str("Group name", request.Distribute.GroupName).Msg("Request distribution group by name")
		group, err := s.requestGroup(request.Distribute.GroupName, request.OwnerName, request.AppName)
		if err != nil {
			return err
		}

		log.Info().Str("Group ID", group.ID).Msg("Distributing to group")
		err = s.releaseToGroup(ctx, request.OwnerName, request.AppName, releaseID, group.ID)
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}

func (s *DistributeService) requestGroup(groupName string,
	ownerName string,
	appName string) (*distributionGroupResponse, error) {

	log.Info().Str("Group name", groupName).Msg("Requesting group")

	// Create Request
	req, err := http.NewRequest("GET",
		fmt.Sprintf("%s/apps/%s/%s/distribution_groups/%s",
			s.client.BaseURL,
			ownerName,
			appName,
			groupName), nil)

	req = s.client.ApplyTokenToRequest(req)

	if err != nil {
		return nil, err
	}

	// Do the request
	response := &distributionGroupResponse{}
	resp, err := s.client.do(req, response)
	if err != nil {
		return nil, err
	}

	if resp.StatusError != nil {
		return nil, resp.StatusError
	}

	log.Info().Str("Group ID", response.ID).Msg("Group ID resolved successfully")

	return response, nil
}

func (s *DistributeService) releaseToGroup(
	ctx context.Context,
	ownerName string,
	appName string,
	releaseID int64,
	groupID string) error {

	body := distributionBody{
		ID:              groupID,
		MandatoryUpdate: false,
		NotifyTester:    false,
	}

	r := distributionResponse{}

	path := fmt.Sprintf("releases/%v/groups", releaseID)

	err := s.client.NewAPIRequest(ctx, http.MethodPost, path, &body, &r)
	if err != nil {
		return err
	}

	return nil
}
