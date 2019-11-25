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

First, install Go1.13+, Node, and Python3 for your environment. For MacOS,

```
brew install go node python3
```

Then, install some libraries and tools:
```
go get -u github.com/magefile/mage golang.org/x/lint/golint golang.org/x/tools/cmd/goimports
pip3 install --upgrade awscli cfn-lint
```

Now you can run `mage` to see the list of available commands.
For example, to run all required checks: `mage test:ci` (`-v` for verbose mode).

## Repo Structure
Since the majority of Panther is written in Go, we follow the [standard Go project layout](https://github.com/golang-standards/project-layout):

* [`docs/`](docs) - code owners, contributing guidelines, etc
* [`internal/`](internal) - backend components
    * [`compliance/`](internal/compliance) - infrastructure scanning and compliance evaluation
    * [`logproc/`](internal/logproc) - security log processing and analysis
    * [`shared/`](internal/shared) - backend components applicable to both products
* [`pkg/`](pkg) - shared standalone packages that could also be imported by other projects
* [`web/`](web) - web application source
