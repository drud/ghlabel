# ghlabel
[![CircleCI](https://circleci.com/gh/drud/ghlabel.svg?style=shield)](https://circleci.com/gh/drud/ghlabel) [![Go Report Card](https://goreportcard.com/badge/github.com/drud/ghlabel)](https://goreportcard.com/report/github.com/drud/ghlabel) ![project is maintained](https://img.shields.io/maintenance/yes/2017.svg)

Manually creating issue labels for each repository in an organization is a tedious and error prone process that simply doesn't scale.

Ghlabel aims to solve this problem by being a tool to automatically standardize GitHub issue labels.
It does so by using issue labels from a reference repository, and applying those labels to all or
a single repository in a GitHub organization or user's account.

## Quickstart
Before getting started, you need to have an API token from GitHub to access any repositories. If you don't already have a token, you can get one [here](https://github.com/settings/tokens).

After the environment variable for the GitHub token is set, you're ready to go.

### Download and install

#### Official release
We recommend downloading ghlabel using the latest release which is available [here](https://github.com/drud/ghlabel/releases).

### Manual install
For ghlabel, we use `make` to generate the executables for all operating systems and architectures.
```
$ export GHLABEL_GITHUB_TOKEN=1234...
$ make
```

## Usage
The tool currently has two functions: previewing staged label changes and applying them.

As a safeguard, ghlabel runs in preview mode by default.
```
./ghlabel --org=drud --ref=community
```
You can apply label changes using the `--apply` flag (or `-a` for short).
```
./ghlabel --org=drud --ref=community --apply
```