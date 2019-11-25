package main

import (
	"fmt"
	"goappcenter/appcenter"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {

	app := cli.App{
		Name: "go-appcenter",
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "apiKey",
			Required: true,
			Usage:    "AppCenter.ms API key",
		},
		&cli.StringFlag{
			Name:     "appName",
			Required: true,
			Usage:    "AppCenter app name",
		},
		&cli.StringFlag{
			Name:     "ownerName",
			Required: true,
			Usage:    "AppCenter owner name",
		},
		&cli.StringFlag{
			Name:     "buildNumber",
			Required: false,
			Usage:    "Release build number",
		},
		&cli.StringFlag{
			Name:     "buildVersion",
			Required: false,
			Usage:    "Release build version",
		},
		&cli.IntFlag{
			Name:     "releaseId",
			Required: false,
			Usage:    "Release version Id",
		},
	}
	app.Name = "Golang AppCenter.ms"
	app.Usage = "Upload and distribute binaries on the AppCenter platform"
	app.Commands = []*cli.Command{
		&cli.Command{
			Name: "upload",
			Flags: []cli.Flag{
				&cli.PathFlag{Name: "file",
					Aliases:  []string{"f"},
					Required: true},
			},
			Action: executeUpload,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println("\t", err)
	}
}

func executeUpload(c *cli.Context) error {
	o := appcenter.ReleaseUploadPayload{
		BuildNumber:  c.String("buildNumber"),
		BuildVersion: c.String("buildVersion"),
		ReleaseID:    c.Int("releaseId"),
	}

	r := appcenter.UploadRequest{
		OwnerName: c.String("ownerName"),
		AppName:   c.String("appName"),
		FilePath:  c.Path("file"),
		Option:    o,
	}
	return appcenter.NewClient(c.String("apiKey")).Upload.Do(r)
}
