![Panther Logo](docs/img/logo-banner.png)

[![Built with Mage](https://magefile.org/badge.svg)](https://magefile.org)
[![CircleCI](https://circleci.com/gh/panther-labs/panther.svg?style=svg)](https://circleci.com/gh/panther-labs/panther)

---

Panther is an open source, cloud-native SIEM written in Golang/React.

Developed by a [dedicated team](https://runpanther.io/about/) of cloud security practitioners, Panther is designed to be:

- **Flexible:** Python-based detections with integrations into common tools such as PagerDuty, Slack, MS Teams, and more.
- **Scalable:** Built on serverless for cost and operational efficiency at any scale.
- **Secure:** Least-privilege and encrypted infrastructure that you control.
- **Integrated:** Support for many popular security logs used for incident response, combined with rich information about your cloud resources.
- **Automated:** Quick and simple deployments with AWS CloudFormation.

Panther use cases:

- **SIEM:** Centralize all security log data for detection, long-term storage, and investigations.
- **[Threat Detection](https://runpanther.io/log-analysis):** Detect suspicious activity quickly and effectively with powerful Python rules.
- **Alerting:** Send notifications to your team when new issues are identified.
- **[Cloud Compliance](https://runpanther.io/compliance/):** Ensure AWS infrastructure abides by defined Python policies.
- **Automatic Remediation:** Fix insecure infrastructure in any number of accounts.

Check out our [website](https://runpanther.io), [blog](https://blog.runpanther.io), and [docs](https://docs.runpanther.io) to learn more.

> **_NOTE:_** Panther is under active development and may experience breaking changes.

## Setup

Panther depends on you having Docker 1.17+ installed. Please refer to the [official Docker page](https://docs.docker.com/install/)
for information on how to install it on your OS.

Having Docker installed, the next step would be to install the necessary build tools. You have the option to either install
them manually on your machine or use the [docker image](https://hub.docker.com/r/pantherlabs/panther-buildpack) that we provide.

### Manually

Install Go1.13+, Node10+, Python3.7+, [Mage](https://magefile.org/#installation), and the [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv1.html)

```bash
brew install go node python3  # MacOS

go get -u -d github.com/magefile/mage
cd $GOPATH/src/github.com/magefile/mage
go run bootstrap.go
```

Make sure that `$GOPATH` is setup correctly and then run `mage setup` to install the prerequisite development libraries.

### With Docker

On the project's root directory, run:

```bash
docker run -v $(pwd):/code -it pantherlabs/panther-buildpack bash # add `-m=4gb` flag if more memory is needed
cd code
```

Run `mage setup` to install the prerequisite development libraries.

## Workflows

Run `mage` to see the list of available commands (`-v` for verbose mode).

You can easily chain `mage` commands together, for example:

```bash
mage fmt test:ci deploy test:integration
```

### Develop

Typical developer workflows.

```bash
mage build:api  # generate Go SDKs for Panther APIs
mage fmt        # format all code
mage test:ci    # run all required checks
```

### Deploy

If you haven't already, [configure](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html) your AWS region and credentials.

Then deploying is as simple as `mage deploy`! You will be prompted to enter a name and email for
the default admin user. Once the deploy is complete, that email will receive a link to sign in.

> NOTE (1): The initial deploy will take 10-15 minutes, and the `deploy` command may timeout before the stack is
> actually finished. Check the AWS CloudFormation console for the status of your deployment.

> NOTE (2): The emails are sent via **no-reply@verificationemail.com**. If you don't see them in
> your inbox, be sure to check the spam folder as well.

# Testing Panther

We strongly recommend running our tests locally before you submit any PR.

```bash
mage test:integration  # Run all integration tests
PKG=./internal/compliance/compliance-api/main mage test:integration  # Run tests for only one package
```

To facilitate testing, we maintain [an image](https://hub.docker.com/r/pantherlabs/panther-buildpack) with all the necessary tools needed to run your tests.
The `Dockerfile` for this image can be found at [build/ci/Dockerfile](build/ci/Dockerfile).

In order to use this, you have to have [Docker](https://www.docker.com/) installed and just spin up and interactive container
with our `pantherlabs/panther-buildpack` image:

```
docker container run -it pantherlabs/panther-buildpack bash
```

From there, you can run any of the testing commands that we provide

## Maintaining the Testing Image

Should a problem arise with
the testing environment (e.g. a version of a certain tool is too old), feel free to contribute
by modifying the `Dockerfile` accordingly.

In order to publish the Dockerfile changes to our image repo (supposing that you are already connected to docker hub
and that you are a member of `pantherlabs`), you have to:

#### 1. Build & Tag the image

To build the image from the modified `Dockerfile` and tag it for a release, run the following from the project's root directory:

```
docker build
    -t pantherlabs/panther-buildpack:<NEW_VERSION> \
    -t pantherlabs/panther-buildpack:latest \
    -f ./build/ci/Dockerfile \
    .
```

In order to decide what the `<NEW_VERSION>` should be, inspect the [latest published image](https://hub.docker.com/r/pantherlabs/panther-buildpack)
and decide what its new value should be according to [semver rules](https://semver.org/).

#### 2. Push the newly built & tagged image to docker hub

To finally push the image, run:

```
docker push pantherlabs/panther-buildpack
```

## CI Service

We utilise a public CircleCI project to run our test suite. The configuration can be found under `.circleci/config.yml`
and all of our jobs use our custom panther docker image. If you create a deploy a new image, don't forget to also modify
the CircleCI executor configuration.

The CircleCI project can be found [here](https://circleci.com/gh/panther-labs/panther/)

## Repo Structure

Since the majority of Panther is written in Go, we follow the [standard Go project layout](https://github.com/golang-standards/project-layout):

- [`api/`](api) - Go models and Swagger specs for communicating with Panther APIs
- [`deployments/`](deployments) - CloudFormation templates for all of Panther's infrastructure
- [`docs/`](docs) - code owners, contributing guidelines, etc
- [`internal/`](internal) - frontend & backend components
- [`pkg/`](pkg) - shared standalone packages that could also be imported by other projects
- [`tools/`](tools) - source code for mage targets
- [`web/`](web) - web application source

## Contributing

Please read the [contributing guidelines](https://github.com/panther-labs/panther/blob/master/docs/CONTRIBUTING.md) before submitting pull requests.
