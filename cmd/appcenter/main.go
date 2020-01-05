package main

import (
	"goappcenter/appcenter"
	"os"

	"github.com/urfave/cli/v2"
)

// AppCenter API Key
var APIKey string

var request = appcenter.UploadRequest{
	Distribute: appcenter.DistributionPayload{},
	Option:     appcenter.ReleaseUploadPayload{},
}

func main() {

	app := cli.App{
		Name:    "go-appcenter",
		Version: "0.1.0",
	}

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Destination: &APIKey,
			EnvVars:     []string{"AppCenterAPIKey"},
			Name:        "apiKey",
			Required:    true,
			Usage:       "AppCenter.ms API key",
		},
	}
	app.Name = "Golang AppCenter.ms"
	app.Usage = "Upload and distribute binaries on the AppCenter platform"
	app.Commands = []*cli.Command{
		&cli.Command{
			Name: "upload",
			Flags: []cli.Flag{
				&cli.PathFlag{Name: "file",
					EnvVars:     []string{"AppCenterFileName"},
					Aliases:     []string{"f"},
					Destination: &request.FilePath,
					Required:    true,
				},
				&cli.StringFlag{
					EnvVars:     []string{"AppCenterAppName"},
					Destination: &request.AppName,
					Name:        "appName",
					Required:    true,
					Usage:       "AppCenter app name",
				},
				&cli.StringFlag{
					Destination: &request.OwnerName,
					EnvVars:     []string{"AppCenterOwnerName"},
					Name:        "ownerName",
					Required:    true,
					Usage:       "AppCenter owner name",
				},
				&cli.StringFlag{
					Destination: &request.Option.BuildNumber,
					Name:        "buildNumber",
					Required:    false,
					Usage:       "Release build number",
				},
				&cli.StringFlag{
					Destination: &request.Option.BuildVersion,
					Name:        "buildVersion",
					Required:    false,
					Usage:       "Release build version",
				},
				&cli.IntFlag{
					Destination: &request.Option.ReleaseID,
					Name:        "releaseId",
					Required:    false,
					Usage:       "Release version Id",
				},
				&cli.StringFlag{
					Destination: &request.Distribute.GroupName,
					EnvVars:     []string{"groupName"},
					Name:        "groupName",
					Required:    false,
					Usage:       "Group name to distribute to the release",
				},
			},
			Action: executeUpload,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func executeUpload(c *cli.Context) error {
	client := appcenter.NewClient(APIKey)

	releaseID, err := client.Upload.Do(request)
	if err != nil {
		return err
	}

	if request.Distribute.GroupName != "" {
		return client.Distribute.Do(*releaseID, request)
	}

	return nil
}
