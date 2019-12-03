![Panther Logo](docs/img/logo-banner.png)

[![Built with Mage](https://magefile.org/badge.svg)](https://magefile.org)

> **NOTE**: Panther is under active development and is still in alpha - breaking changes are likely

Panther is an open source, cloud-native SIEM written in Golang/React.

Developed by a [dedicated team](https://runpanther.io/about/) of cloud and security experts, Panther is designed to be:

* **Flexible:** Python-based analysis for customized detection
* **Scalable:** Built on a modern, serverless platform
* **Secure:** Least-privilege access to encrypted infrastructure you control
* **Integrated:** Enrich log analysis with information about your cloud

Check out our [website](https://runpanther.io), [blog](https://blog.runpanther.io), and [documentation](https://docs.runpanther.io) to learn more!

## Products
Panther provides two main features: [cloud security](https://runpanther.io/compliance/) and
[threat detection](https://runpanther.io/log-analysis), and provides flexibility to select only the features you need.

## Setup
Install Go1.13+, [Mage](https://magefile.org/#installation), Node, and Python3:

```
brew install go node python3  # MacOS

go get -u -d github.com/magefile/mage
cd $GOPATH/src/github.com/magefile/mage
go run bootstrap.go
```

Then run `mage setup` to install the prerequisite development libraries.
(`mage build` and `mage test` will automatically install Go packages as they are needed.)

## Development and Deployment
Run `mage` to see the list of available commands (`-v` for verbose mode).
Steps in a typical developer workflow might be:

```bash
mage build:api  # generate Go SDKs for Panther APIs
mage fmt        # format all code
mage test:ci    # run all required checks

# Deploy and run integration tests
mage deploy:pre  # deploy prerequisite S3 buckets
AWS_REGION=us-west-2 PARAMS="Debug=true" mage deploy:backend
PKG=./internal/compliance/compliance-api/main mage test:integration
```

## Repo Structure
Since the majority of Panther is written in Go, we follow the [standard Go project layout](https://github.com/golang-standards/project-layout):

* [`api/`](api) - Go models and Swagger specs for communicating with Panther APIs
* [`deployments/`](deployments) - CloudFormation templates for all of Panther's infrastructure
* [`docs/`](docs) - code owners, contributing guidelines, etc
* [`internal/`](internal) - backend components
* [`pkg/`](pkg) - shared standalone packages that could also be imported by other projects
* [`tools/`](tools) - source code for mage targets
* [`web/`](web) - web application source
