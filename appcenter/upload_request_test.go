package appcenter

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValdationBuildVersionArgument(t *testing.T) {
	r := UploadTask{}

	t.Run("When the build verson number is missing", func(t *testing.T) {
		testCases := []struct {
			ext          string
			buildVersion string
			buildNumber  string
			err          bool
		}{
			{"zip", "", "", true},
			{"zip", "1.2.3", "", false},
			{"zip", "", "1", true},

			{"msi", "", "", true},
			{"msi", "1.2.3", "", false},
			{"msi", "", "1", true},

			{"apk", "", "", false},
			{"apk", "1.2.3", "", false},
			{"apk", "1.2.3", "1", false},
			{"apk", "", "1", false},

			{"pkg", "", "", true},
			{"pkg", "1.2.3", "", true},
			{"pkg", "", "1", true},
			{"pkg", "1.2.3", "1", false},

			{"dmg", "", "", true},
			{"dmg", "1.2.3", "", true},
			{"dmg", "", "1", true},
			{"dmg", "1.2.3", "1", false},
		}

		for _, tc := range testCases {
			t.Run(fmt.Sprintf("For ext: %v BuildVersion: %v BuildNumber: %v",
				tc.buildVersion, tc.buildNumber, tc.ext), func(t *testing.T) {

				r.Option.BuildVersion = tc.buildVersion
				r.Option.BuildNumber = tc.buildNumber

				r.FilePath = fmt.Sprintf("toto.%v", tc.ext)
				err := r.validateRequestBuildVersion()
				if tc.err {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
}
