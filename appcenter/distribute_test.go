package appcenter

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mc *MockClient
var ds DistributeService
var req http.Request

func setup() {
	mc = new(MockClient)
	req = http.Request{}
	mc.On("ApplyTokenToRequest", &req.Header).Return()
	ds = DistributeService{client: mc}
}

func TestRequestGroup(t *testing.T) {
	t.Run("Should handle service request errors", func(t *testing.T) {
		// setup:
		setup()

		mc.On("NewServiceRequest",
			"GET",
			"/apps/toto/appname/distribution_groups/Group name",
			mock.Anything).Return(&req, errors.New("Test error"))

		// when:
		groupID, err := ds.requestGroup("Group name", "toto", "appname")

		// then:
		assert.EqualError(t, err, "Test error")

		// and:
		assert.Empty(t, groupID)
	})

	t.Run("Should handle do error", func(t *testing.T) {
		// setup:
		setup()

		mc.On("NewServiceRequest",
			"GET",
			"/apps/toto/appname/distribution_groups/Group name",
			mock.Anything).
			Return(&req, nil)

		mc.On("Do",
			mock.Anything,
			mock.Anything).
			Return(&Response{}, errors.New("Test error"))

		// when:
		groupID, err := ds.requestGroup("Group name", "toto", "appname")
		assert.Empty(t, groupID)
		assert.EqualError(t, err, "Test error")
	})

	t.Run("Should parse ID", func(t *testing.T) {
		// setup:
		setup()

		mc.On("NewServiceRequest",
			"GET",
			"/apps/toto/appname/distribution_groups/Group name",
			mock.Anything).
			Return(&req, nil)

		mc.On("Do",
			mock.Anything,
			mock.Anything).
			Return(&Response{}, nil)

		// when:
		groupID, err := ds.requestGroup("Group name", "toto", "appname")
		fmt.Println(groupID, err)
		assert.Empty(t, groupID)
		assert.NoError(t, err)
	})
}
