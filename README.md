# go-appcenter

![dl](https://img.shields.io/docker/pulls/sho3box/go-appcenter.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/shoebox/go-appcenter)](https://goreportcard.com/report/github.com/shoebox/go-appcenter)
[![Build
Status](https://travis-ci.org/shoebox/go-appcenter.svg)](https://travis-ci.org/shoebox/go-appcenter)

`AppCenter.ms` upload and distribution client made in Go language.

# Features

- ✅ Upload and (optionally) distribute a binary to the AppCenter.ms platform
- ✅ Supports APK, IPA, PKG, DMG, ZIP, MSI... upload
- ✅ Statically compiled, do not require any runtime dependencies
- ✅ Parallelized chunks upload
- ✅ Up-to-date with latest API

# Usage

## Demo

![demo](demo.gif)

## Command line interface

Basic help:

```bash
NAME:
   Golang AppCenter.ms - Upload and distribute binaries on the AppCenter platform

USAGE:
   main [global options] command [command options] [arguments...]

VERSION:
   0.2.0

COMMANDS:
   upload
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --apiKey value  AppCenter.ms API key [$AppCenterAPIKey]
   --help, -h      show help (default: false)
   --version, -v   print the version (default: false)
```

## Upload command

### Arguments

| Arg              | Mandatory | Description                                                                                                    |
| ---              | ---       | ---                                                                                                            |
| `--file`         | YES       | [AppCenter API Key](https://docs.microsoft.com/en-us/appcenter/api-docs/#creating-an-app-center-app-api-token) |
| `--appName`      | YES       | Application name in AppCenter                                                                                  |
| `--ownerName`    | YES       | Application owner in AppCenter                                                                                 |
| `--buildNumber`  | NO        | Build number                                                                                                   |
| `--buildVersion` | NO        | Build version string                                                                                           |
| `--releaseId`    | NO        | Release ID                                                                                                     |
| `--groupName`    | NO        | Group to distribute binary to                                                                                  |

### Arguments as environment values

Command arguments can be configured via environment variables.

The variable in question being:

| Name               | Description                 | 
| ---                | ---               |
| AppCenterAPIKey    | AppCenter API Key           |
| AppCenterOwnerName | AppCenter application owner | 
| AppCenterAppName   | AppCenter application name  |


### How resolve AppName and OwnerName in AppCenter

Refer to the application URL in AppCenter:

`https://appcenter.ms/orgs/<OWNER_NAME>/apps/<APP_NAME>`

### Help

```bash
NAME:
   main upload -

USAGE:
   main upload [command options] [arguments...]

DESCRIPTION:
   Upload binary to AppCenter for distribution. And optionally distribute it

OPTIONS:
   --file value, -f value   [$AppCenterFileName]
   --appName value         AppCenter app name [$AppCenterAppName]
   --ownerName value       AppCenter owner name [$AppCenterOwnerName]
   --buildNumber value     Release build number
   --buildVersion value    Release build version
   --releaseId value       Release version Id (default: 0)
   --groupName value       Group name to distribute to the release [$groupName]
   --help, -h              show help (default: false)
```

## Via Docker

Image is hosted on [DockerHub](https://hub.docker.com/r/sho3box/go-appcenter)

To run it:

```bash
docker run sho3box/go-appcenter:latest
```


# Dependencies

The library uses a GO module to manage dependencies.

# License

go-appcenter is licensed under the MIT license. See [LICENSE](LICENSE) for more info.
