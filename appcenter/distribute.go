package appcenter

import (
	"fmt"

	"github.com/fatih/color"
)

// DistributeService definition
type DistributeService struct {
	client AppCenterClient
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

// Do Distribute the designated release into the provided configuration
func (s *DistributeService) Do(releaseID string, request UploadRequest) error {
	if request.Distribute.GroupName != "" {
		color.Green("\n\tDistributing release")
		groupID, err := s.requestGroup(request.Distribute.GroupName, request.OwnerName, request.AppName)
		if err != nil {
			return err
		}

		err = s.releaseToGroup(request.OwnerName, request.AppName, releaseID, groupID)
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}

func (s *DistributeService) requestGroup(groupName string,
	ownerName string,
	appName string) (string, error) {

	fmt.Println("\t\tRequesting group", groupName)

	// Create Request
	path := s.computePath("distribution_groups", ownerName, appName, groupName)

	response := &distributionGroupResponse{}
	req, err := s.client.NewServiceRequest("GET", path, response)
	s.client.ApplyTokenToRequest(&req.Header)

	if err != nil {
		return "", err
	}

	// Do the request
	_, err = s.client.Do(req, response)
	if err != nil {
		return "", err
	}

	fmt.Println("\t\tGroup ID :", response.ID)

	return response.ID, nil
}

func (s *DistributeService) computePath(serviceName string, ownerName string, appName string, groupName string) string {
	return fmt.Sprintf("/apps/%s/%s/%s/%s",
		ownerName,
		appName,
		serviceName,
		groupName)
}

func (s *DistributeService) releaseToGroup(ownerName string,
	appName string,
	releaseID string,
	groupID string) error {

	url := fmt.Sprintf("/apps/%s/%s/releases/%s/groups",
		ownerName,
		appName,
		releaseID)

	body := distributionBody{
		ID:              groupID,
		MandatoryUpdate: false,
		NotifyTester:    false,
	}

	req, err := s.client.NewServiceRequest("POST", url, body)
	if err != nil {
		return err
	}

	s.client.ApplyTokenToRequest(&req.Header)
	s.client.RequestContentTypeJSON(&req.Header)

	resp, err := s.client.Do(req, body)
	if err != nil {
		return err
	}

	if resp.Response.StatusCode == 201 {
		color.Green("\tDistribution completed")
		return nil
	}

	// TODO: Wrap better the error here
	return fmt.Errorf("Failed to share release %v to group %v (Error : %v)",
		releaseID,
		groupID,
		resp)
}
