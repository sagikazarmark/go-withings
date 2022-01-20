# Go client for the [Withings API](https://developer.withings.com/)

[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/sagikazarmark/go-withings/CI?style=flat-square)](https://github.com/sagikazarmark/go-withings/actions?query=workflow%3ACI)
[![Codecov](https://img.shields.io/codecov/c/github/sagikazarmark/go-withings?style=flat-square)](https://codecov.io/gh/sagikazarmark/go-withings)
[![Go Report Card](https://goreportcard.com/badge/github.com/sagikazarmark/go-withings?style=flat-square)](https://goreportcard.com/report/github.com/sagikazarmark/go-withings)
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.16-61CFDD.svg?style=flat-square)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/mod/github.com/sagikazarmark/go-withings)
[![built with nix](https://img.shields.io/badge/builtwith-nix-7d81f7?style=flat-square)](https://builtwithnix.org)

**go-withings is a Go client library for accessing the [Withings API](https://developer.withings.com/).**

**⚠️ WARNING: This is still work in progress. ⚠️**


## Installation

```shell
go get github.com/sagikazarmark/go-withings
```

## API coverage

The Withings API provides a wide range of services, but many of them are targeted at (health) service providers.
The primary focus of this SDK is to provide access to the data APIs, so providing a full coverage is not a goal at this time.
That being said, PRs are always welcome.

Supported API services/calls:

- OAuth2
- Measure (WIP)
- Heart (WIP)
- Sleep (WIP)
- Notify (WIP)

Unsupported API services/calls:

- Dropshipment
- User
- Signature

Feel free to open a discussion or issue if something is missing and you would like it to be included.


## Development

When all coding and testing is done, please run the test suite:

```shell
make check
```

For the best developer experience, install [Nix](https://builtwithnix.org/) and [direnv](https://direnv.net/).

Alternatively, install Go manually or using a package manager. Install the rest of the dependencies by running `make deps`.
