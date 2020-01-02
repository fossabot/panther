![Panther Logo](docs/img/logo-banner.png)

[![Built with Mage](https://magefile.org/badge.svg)](https://magefile.org)

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

Install Go1.13+, Node, Python3, [Mage](https://magefile.org/#installation), and the [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv1.html)

```bash
brew install go node python3  # MacOS

go get -u -d github.com/magefile/mage
cd $GOPATH/src/github.com/magefile/mage
go run bootstrap.go
```

Then run `mage setup` to install the prerequisite development libraries.

Finally, configure the required fields in [panther_config.yml](deployments/panther_config.yml).

## Workflows

Run `mage` to see the list of available commands (`-v` for verbose mode).

You can easily chain `mage` commands together, for example:

```bash
mage fmt test:ci deploy:app test:integration
```

### Develop

Typical developer workflows.

```bash
mage build:api  # generate Go SDKs for Panther APIs
mage fmt        # format all code
mage test:ci    # run all required checks
```

### Deploy

Before you deploy, please make sure you fill in the email of the user that we are going
to automatically create for you. Please edit `deployments/panther_config.yml` and fill in the
`UserEmail` parameter with the email of administrator in your team.

After that

```bash
mage deploy:pre  # deploy prerequisite S3 buckets (only need to do once)
mage deploy:app

# Optional: Deploy with parameters
export AWS_REGION=us-west-2
export PARAMS="Debug=true"
mage deploy:app
```

### Integration Testing

Run tests on the deployed infrastructure to ensure each component is operating as intended.

```bash
mage test:integration  # Run all integration tests
PKG=./internal/compliance/compliance-api/main mage test:integration  # Run tests for only one package
Creates the necessary AWS infrastructure to setup the main Panther application.
```

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
