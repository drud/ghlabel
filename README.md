# ghlabel
A CLT to manage GitHub labels in accordance with [community guidelines](https://github.com/drud/community/blob/master/development/issue_workflow.md#labels).

## Setup
Clone this repository and `cd` into the project. Once you're in the project's root folder, run `make <your_sys_arch>`.
Set the $GITHUB_TOKEN environment variable using a GitHub API Key.
```
export GITHUB_TOKEN=12345...
```

## Usage
The tool currently allows you to preview and execute proposed label changes.

Using preview mode allows you to see label changes before they are executed. It's called like this:
```
./ghlabel --owner=drud --parent=community preview
```
Intuitively, you execute changes staged in preview mode by calling `ghlabel` with no arguments.
```
./ghlabel --owner=drud --parent=community
```