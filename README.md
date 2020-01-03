# go-appcenter


[![codecov](https://codecov.io/gh/shoebox/go-appcenter/branch/master/graph/badge.svg)](https://codecov.io/gh/shoebox/go-appcenter)
[![Go Report Card](https://goreportcard.com/badge/github.com/shoebox/go-appcenter)](https://goreportcard.com/report/github.com/shoebox/go-appcenter)
[![Build
Status](https://travis-ci.org/shoebox/go-appcenter.svg)](https://travis-ci.org/shoebox/go-appcenter)

AppCenter.ms client made in GO

# Features

- ✅ Upload and distribute binaries to the AppCenter.ms platform
- ✅ Supports APK, IPA, PKG, DMG, ZIP, MSI upload

# Usage

## Via Docker

Image is hosted on DockerHub:
> https://hub.docker.com/repository/docker/sho3box/go-appcenter

To run it:

> docker run sho3box/go-appcenter:latest

## Via GO CLI

	go run .

Will return:

	NAME:
	Golang AppCenter.ms - Upload and distribute binaries on the AppCenter platform

	USAGE:
	go-appcenter [global options] command [command options] [arguments...]

	VERSION:
	0.0.0

	COMMANDS:
	upload
	help, h  Shows a list of commands or help for one command

	GLOBAL OPTIONS:
	--apiKey value  AppCenter.ms API key [$AppCenterAPIKey]
	--help, -h      show help (default: false)
	--version, -v   print the version (default: false)

# Dependencies

go-appcenter uses go module to manage dependencies.

# License

go-appcenter is licensed under the MIT license. See [LICENSE](LICENSE) for more info.
