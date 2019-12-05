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

TBD

## Via GO CLI

	go run .

Will return:

	NAME:
		Golang AppCenter.ms - Upload and distribute binaries on the AppCenter platform

	USAGE:
		goappcenter [global options] command [command options] [arguments...]

	VERSION:
		0.0.0

	COMMANDS:
		upload
		help, h  Shows a list of commands or help for one command

	GLOBAL OPTIONS:
		--apiKey value        AppCenter.ms API key
		--appName value       AppCenter app name
		--ownerName value     AppCenter owner name
		--buildNumber value   Release build number
		--buildVersion value  Release build version
		--releaseId value     Release version Id (default: 0)
		--help, -h            show help (default: false)
		--version, -v         print the version (default: false)
			errr :  Required flags "apiKey, appName, ownerName" not set

# Dependencies

go-appcenter uses go module to manage dependencies.

# License

go-appcenter is licensed under the MIT license. See [LICENSE](LICENSE) for more info.
