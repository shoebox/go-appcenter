package appcenter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fatih/color"
)

// DistributeService definition
type DistributeService struct {
	client *Client
}

type distributionGroupResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Origin      string `json:"origin"`
	displayName string `json:"display_name"`
}

type distributionBody struct {
	ID              string `json:"id"`
	MandatoryUpdate bool   `json:"mandatory_update"`
	NotifyTester    bool   `json:"notify_testers"`
}

// Do Distribute the designated release into the provided configuration
func (s *DistributeService) Do(releaseID string, request UploadRequest) error {
	if request.Distribute.GroupName != "" {
		color.Green("\n\tDistributing release")
		group, err := s.requestGroup(request.Distribute.GroupName, request.OwnerName, request.AppName)
		if err != nil {
			return err
		}

		err = s.releaseToGroup(request.OwnerName, request.AppName, releaseID, group.ID)
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

	fmt.Println("\t\tRequesting group", groupName)

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
	_, err = s.client.do(req, response)
	if err != nil {
		return nil, err
	}

	fmt.Println("\t\tGroup ID :", response.ID)

	return response, nil
}

func (s *DistributeService) releaseToGroup(ownerName string,
	appName string,
	releaseID string,
	groupID string) error {

	body := distributionBody{
		ID:              groupID,
		MandatoryUpdate: false,
		NotifyTester:    false,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/apps/%s/%s/releases/%s/groups",
			s.client.BaseURL,
			ownerName,
			appName,
			releaseID),
		bytes.NewBuffer(payload))

	req = RequestContentTypeJson(req)
	req = s.client.ApplyTokenToRequest(req)
	if err != nil {
		return err
	}

	resp, err := s.client.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == 201 {
		color.Green("\tDistribution completed")
		return nil
	} else {
		// TODO: Wrap better the error here
		return fmt.Errorf("Failed to share release %v to group %v (Error : %v)",
			releaseID,
			groupID,
			resp)
	}
	return nil
}
