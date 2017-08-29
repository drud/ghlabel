# ghlabel
Since all company projects must abide by the [Drud community guidelines](https://github.com/drud/community/blob/master/development/issue_workflow.md#labels),
it makes sense to automate some processes. Hence, Ghlabel.

Ghlabel is a tool that automatically standardizes GitHub issue labels across a user or organization's repositories.
A reference repository is used as the template for labels, and those labels are automatically copied to all or
a single repository.

## Quickstart
Before getting started, you need to have an API token from GitHub to access any repositories. If you don't already have
a token, you can get one [here](https://github.com/settings/tokens).

After the environment variable for the GitHub token is set, you're ready to generate your system's binary. For ghlabel, we use
`make` to generate the executables for all operating systems and architectures.
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