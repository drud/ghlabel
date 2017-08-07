# ghlabel
A CLT to manage GitHub labels in accordance with community guidelines.
## Setup
Clone this repository and `cd` into the project. Once you're in the project's root folder, run `make <your_sys_arch>`.
Set the $GITHUB_TOKEN environment variable using a GitHub API Key.
`export GITHUB_TOKEN=12345...`
The tool is currently executed using `./ghlabel --owner=drud --parent community`.
