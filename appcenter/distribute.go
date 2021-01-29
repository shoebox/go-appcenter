package appcenter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pterm/pterm"
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
	GroupID               string `json:"id"`
	MandatoryUpdate       bool   `json:"mandatory_update"`
	ProvisioningStatusURL string `json:"provisioning_status_url"`
}

// Do Distribute the designated release into the provided configuration
func (s *DistributeService) Do(ctx context.Context, releaseID int64, request UploadTask) error {

	if request.Distribute.GroupName != "" {
		group, err := s.requestGroup(ctx, request.Distribute.GroupName, request.OwnerName, request.AppName)
		if err != nil {
			return err
		}

		err = s.releaseToGroup(ctx, request.OwnerName, request.AppName, releaseID, group.ID)
		return err
	}

	return nil
}

func (s *DistributeService) requestGroup(
	ctx context.Context,
	groupName string,
	ownerName string,
	appName string,
) (*distributionGroupResponse, error) {
	var res distributionGroupResponse

	sp, err := pterm.DefaultSpinner.Start(fmt.Sprintf("Requesting distribution group ID from name '%v'", groupName))
	if err != nil {
		return &res, err
	}

	err = s.client.NewAPIRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("distribution_groups/%s", groupName),
		nil,
		&res,
	)

	if err == nil {
		sp.UpdateText(fmt.Sprintf("Distribution group ID resolved: %v", res.ID))
		sp.Success()
	} else {
		sp.Fail()
	}

	return &res, err
}

func (s *DistributeService) releaseToGroup(
	ctx context.Context,
	ownerName string,
	appName string,
	releaseID int64,
	groupID string) error {

	sp, err := pterm.DefaultSpinner.Start("Releasing to group")
	if err != nil {
		return err
	}

	body := distributionBody{
		ID:              groupID,
		MandatoryUpdate: false,
		NotifyTester:    false,
	}

	r := distributionResponse{}

	path := fmt.Sprintf("releases/%v/groups", releaseID)

	err = s.client.NewAPIRequest(ctx, http.MethodPost, path, &body, &r)
	if err != nil {
		sp.Fail()
		return err
	}

	sp.Success()

	return nil
}
