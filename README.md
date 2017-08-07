# ghlabel
A CLT to manage GitHub labels in accordance with community guidelines.
## Setup
Clone this repository and `cd` into the project. Once you're in the project's root folder, run `make <your_sys_arch>`.
Set the $GITHUB_TOKEN environment variable using a GitHub API Key.
`export GITHUB_TOKEN=12345...`
## Usage
The tool currently supports two modes; `preview` and `run`.

Using preview mode allows you to see label changes before they are executed. It's called like this:
```
./ghlabel --owner=drud --parent=community preview
```
Intuitively, run mode will execute the changes staged in preview.
```
./ghlabel --owner=drud --parent=community run
```